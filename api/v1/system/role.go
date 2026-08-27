package v1

import (
	"net/http"
	"south-admin-gin/models/system"
systemServices "south-admin-gin/services/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRolePage 获取角色分页列表
func GetRolePage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", c.DefaultQuery("page_size", "10")))
	name := c.Query("name")

	result, err := systemServices.GetRolePage(page, pageSize, name)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "获取角色列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: result,
	})
}

// GetRoleDetail 获取角色详情
func GetRoleDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	role, err := systemServices.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 404,
			Message: "角色不存在",
		})
		return
	}

	// 获取菜单ID列表
	menuIds := make([]int, 0)
	for _, menu := range role.Menus {
		menuIds = append(menuIds, int(menu.ID))
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
			"created_at":  role.CreatedAt,
			"updated_at":  role.UpdatedAt,
			"authorize":   menuIds,
		},
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var dto models.CreateRoleDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	role, err := systemServices.CreateRole(dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "创建角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "创建成功",
		Data: role,
	})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	var dto models.UpdateRoleDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	role, err := systemServices.UpdateRole(id, dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "更新角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "更新成功",
		Data: role,
	})
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	err = systemServices.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "删除角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "删除成功",
	})
}

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	roles, err := systemServices.GetRoleList()
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "获取角色列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: roles,
	})
}

// GetRoleAuthorize 获取角色授权信息
func GetRoleAuthorize(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Query("roleId"))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	result, err := systemServices.GetRoleAuthorize(roleId)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "获取授权信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: result,
	})
}

// SaveRoleAuthorize 保存角色授权
func SaveRoleAuthorize(c *gin.Context) {
	var dto models.AuthorizeRoleDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	err := systemServices.SaveRoleAuthorize(dto.RoleId, dto.MenuIds)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "保存授权失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "保存成功",
	})
}
