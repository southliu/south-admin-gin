package system

import (
	"south-admin-gin/middleware"

	"github.com/gin-gonic/gin"
)

// InitSystemRoutes 初始化所有系统路由
func InitSystemRoutes(r *gin.RouterGroup) {
	InitUserRoutes(r.Group("/user"))
	InitRoleRoutes(r.Group("/role"))
	InitMenuRoutes(r.Group("/menu"))
	InitPermissionRoutes(r.Group("/permission"))
	InitLogRoutes(r.Group("/log"))
}

func withAuth(r *gin.RouterGroup) *gin.RouterGroup {
	auth := r.Group("")
	auth.Use(middleware.JWTAuth())
	return auth
}
