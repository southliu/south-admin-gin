package router

import (
	v1 "south-admin-gin/api/v1"

	"github.com/gin-gonic/gin"
)

// initDashboardRoutes 仪表盘及下拉演示数据路由（数据写死）
func initDashboardRoutes(r *gin.Engine) {
	r.GET("/dashboard", v1.GetDashboardTrends)
	r.GET("/platform/partner", v1.GetPartnerList)
	r.GET("/authority/common/games", v1.GetGameList)
}
