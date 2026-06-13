package system

import (
	v1 "serve-wechat-gin/api/v1/system"

	"github.com/gin-gonic/gin"
)

func InitRoleRoutes(r *gin.RouterGroup) {
	auth := withAuth(r)
	auth.GET("/page", v1.GetRolePage)
	auth.POST("/create", v1.CreateRole)
	auth.PUT("/update/:id", v1.UpdateRole)
	auth.GET("/detail", v1.GetRoleDetail)
	auth.GET("/list", v1.GetRoleList)
	auth.DELETE("/:id", v1.DeleteRole)
	auth.GET("/authorize", v1.GetRoleAuthorize)
	auth.PUT("/authorize/save", v1.SaveRoleAuthorize)
}
