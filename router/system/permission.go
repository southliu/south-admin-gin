package system

import (
	v1 "serve-wechat-gin/api/v1/system"

	"github.com/gin-gonic/gin"
)

func InitPermissionRoutes(r *gin.RouterGroup) {
	auth := withAuth(r)
	auth.GET("/page", v1.GetPermissionPage)
	auth.POST("/create", v1.CreatePermission)
	auth.GET("/detail", v1.GetPermissionDetail)
	auth.PUT("/update/:id", v1.UpdatePermission)
	auth.DELETE("/:id", v1.DeletePermission)
	auth.GET("/list", v1.GetPermissionList)
}
