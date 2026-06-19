package v1

import (
	"net/http"
	"south-admin-gin/models/system"
systemServices "south-admin-gin/services/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPermissionPage 获取权限分页列表
func GetPermissionPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")

	result, err := systemServices.GetPermissionPage(page, pageSize, name)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "获取权限列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: result,
	})
}

// GetPermissionDetail 获取权限详情
func GetPermissionDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	permission, err := systemServices.GetPermissionByID(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 404,
			Message: "权限不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: permission,
	})
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var dto models.CreatePermissionDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	permission, err := systemServices.CreatePermission(dto.Name, dto.Description)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "创建权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "创建成功",
		Data: permission,
	})
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	var dto models.UpdatePermissionDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	permission, err := systemServices.UpdatePermission(id, dto.Name, dto.Description)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "更新权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "更新成功",
		Data: permission,
	})
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 400,
			Message: "参数错误",
		})
		return
	}

	err = systemServices.DeletePermission(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "删除权限失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "删除成功",
	})
}

// GetPermissionList 获取权限列表
func GetPermissionList(c *gin.Context) {
	permissions, err := systemServices.GetPermissionList()
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code: 500,
			Message: "获取权限列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code: 200,
		Message: "获取成功",
		Data: permissions,
	})
}
