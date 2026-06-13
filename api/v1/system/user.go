package v1

import (
	"log"
	"net/http"
	"serve-wechat-gin/models/system"
	systemServices "serve-wechat-gin/services/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
func Login(c *gin.Context) {
	var req models.LoginData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.CommonResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 查询数据库验证用户名密码
	user, err := systemServices.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "用户名或密码错误",
		})
		return
	}

	// 校验密码
	if !systemServices.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "用户已被禁用",
		})
		return
	}

	// 生成 JWT
	token, err := systemServices.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成 token 失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.CommonResponse{
			Code:    500,
			Message: "登录失败，请重试",
		})
		return
	}

	// 获取用户权限
	permissions := systemServices.GetUserPermissions(user)

	// 获取用户角色ID列表
	roles := make([]int, 0)
	for _, role := range user.Roles {
		roles = append(roles, int(role.ID))
	}

	result := models.LoginResult{
		Token: token,
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Phone:    user.Phone,
			Email:    user.Email,
			Status:   user.Status,
			Roles:    roles,
		},
		Roles:       roles,
		Permissions: permissions,
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "登录成功",
		Data:    result,
	})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "退出登录成功",
	})
}

// RefreshPermissions 刷新用户权限
func RefreshPermissions(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    401,
			Message: "未登录",
		})
		return
	}

	result, err := systemServices.RefreshPermissions(userId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "刷新权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	})
}

// GetUserPage 获取用户分页列表
func GetUserPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")

	result, err := systemServices.GetUserPage(page, pageSize, username)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "获取用户列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	})
}

// GetUserDetail 获取用户详情
func GetUserDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	user, err := systemServices.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	// 获取角色ID列表
	roleIds := make([]int, 0)
	for _, role := range user.Roles {
		roleIds = append(roleIds, int(role.ID))
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"status":     user.Status,
			"role_ids":   roleIds,
			"roles":      user.Roles,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var dto models.CreateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	user, err := systemServices.CreateUser(dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "创建用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "创建成功",
		Data:    user,
	})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	var dto models.UpdateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	user, err := systemServices.UpdateUser(id, dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "更新用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "更新成功",
		Data:    user,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	err = systemServices.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "删除用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "删除成功",
	})
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	users, err := systemServices.GetUserList()
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "获取用户列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    users,
	})
}

// UpdatePassword 更新密码
func UpdatePassword(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    401,
			Message: "未登录",
		})
		return
	}

	var dto models.UpdatePasswordDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	err := systemServices.UpdatePassword(userId.(int64), dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "更新密码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "更新成功",
	})
}
