package v1

import (
	"net/http"
	"south-admin-gin/models/system"

	"github.com/gin-gonic/gin"
)

// GetDashboardTrends 仪表盘数据总览（演示数据，前端仅依赖接口成功返回）
func GetDashboardTrends(c *gin.Context) {
	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data: gin.H{
			"payDate":  "2022-10-19 ~ 2022-10-29",
			"user":     9424,
			"recharge": 2590,
			"order":    3478,
			"game":     350,
		},
	})
}

// GetPartnerList 合作公司下拉数据（演示数据）
func GetPartnerList(c *gin.Context) {
	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data: []gin.H{
			{"id": "1", "name": "演示公司A"},
			{"id": "2", "name": "演示公司B"},
			{"id": "3", "name": "演示公司C"},
		},
	})
}

// GetGameList 游戏下拉数据（演示数据）
func GetGameList(c *gin.Context) {
	c.JSON(http.StatusOK, models.CommonResponse{
		Code:    200,
		Message: "获取成功",
		Data: []gin.H{
			{"id": "1", "name": "演示游戏A"},
			{"id": "2", "name": "演示游戏B"},
			{"id": "3", "name": "演示游戏C"},
		},
	})
}
