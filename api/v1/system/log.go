package v1

import (
	"net/http"
	"south-admin-gin/models/system"
systemServices "south-admin-gin/services/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateLog 创建日志
func CreateLog(c *gin.Context) {
	var dto models.CreateLogDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, models.CommonResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := systemServices.CreateLog(dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "创建日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "创建成功",
	})
}

// GetLogPage 获取日志分页列表
func GetLogPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	logType, _ := strconv.Atoi(c.DefaultQuery("type", "-1"))

	result, err := systemServices.GetLogPage(page, pageSize, username, logType)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "获取日志列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data:    result,
	})
}

// DeleteLog 删除日志
func DeleteLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	err = systemServices.DeleteLog(id)
	if err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "删除日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "删除成功",
	})
}

// BatchDeleteLog 批量删除日志
func BatchDeleteLog(c *gin.Context) {
	var dto struct {
		Ids []int64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := systemServices.BatchDeleteLog(dto.Ids); err != nil {
		c.JSON(http.StatusOK, models.CommonResponse{
			Code:    500,
			Message: "批量删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "批量删除成功",
	})
}
