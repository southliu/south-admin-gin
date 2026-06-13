package v1

import (
	"net/http"
	"serve-wechat-gin/models/system"
systemServices "serve-wechat-gin/services/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMenuList 获取用户菜单列表
func GetMenuList(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    401,
			Message: "未登录",
		})
		return
	}

	menus, err := systemServices.GetMenuList(userId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "获取菜单列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    menus,
	})
}

// GetMenuPage 获取菜单分页列表
func GetMenuPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	label := c.Query("label")
	labelEn := c.Query("labelEn")
	rule := c.Query("rule")

	var state *int
	if stateStr := c.Query("state"); stateStr != "" {
		s, _ := strconv.Atoi(stateStr)
		state = &s
	}

	result, err := systemServices.GetMenuPage(page, pageSize, label, labelEn, state, rule)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "获取菜单列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	})
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    401,
			Message: "未登录",
		})
		return
	}

	var dto models.CreateMenuDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	menu, err := systemServices.CreateMenu(dto, userId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "创建菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "创建成功",
		Data:    menu,
	})
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	var dto models.UpdateMenuDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	menu, err := systemServices.UpdateMenu(id, dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "更新菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "更新成功",
		Data:    menu,
	})
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	err = systemServices.DeleteMenu(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "删除菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "删除成功",
	})
}

// GetMenuDetail 获取菜单详情
func GetMenuDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	menu, err := systemServices.GetMenuDetail(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    404,
			Message: "菜单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    menu,
	})
}

// ChangeMenuState 修改菜单状态
func ChangeMenuState(c *gin.Context) {
	var dto models.ChangeMenuStateDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	err := systemServices.ChangeMenuState(dto)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "修改状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "修改成功",
	})
}

// BatchDeleteMenu 批量删除菜单
func BatchDeleteMenu(c *gin.Context) {
	var dto struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if len(dto.IDs) == 0 {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "请选择要删除的菜单",
		})
		return
	}

	err := systemServices.BatchDeleteMenu(dto.IDs)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "批量删除菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "批量删除成功",
	})
}
