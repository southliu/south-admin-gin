package router

import (
	"south-admin-gin/router/system"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	system.InitSystemRoutes(r.Group("/system"))
	return r
}
