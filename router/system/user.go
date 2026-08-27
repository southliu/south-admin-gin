package system

import (
	v1 "south-admin-gin/api/v1/system"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
	r.POST("/login", v1.Login)
	r.POST("/logout", v1.Logout)

	auth := withAuth(r)
	auth.GET("/refreshPermissions", v1.RefreshPermissions)
	auth.GET("/page", v1.GetUserPage)
	auth.GET("/detail", v1.GetUserDetail)
	auth.POST("/create", v1.CreateUser)
	auth.PUT("/update/:id", v1.UpdateUser)
	auth.DELETE("/:id", v1.DeleteUser)
	auth.POST("/batchDelete", v1.BatchDeleteUser)
	auth.GET("/list", v1.GetUserList)
	auth.POST("/updatePassword", v1.UpdatePassword)
}
