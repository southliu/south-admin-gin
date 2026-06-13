package system

import (
	v1 "serve-wechat-gin/api/v1/system"

	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup) {
	auth := withAuth(r)
	auth.GET("/list", v1.GetMenuList)
	auth.GET("/page", v1.GetMenuPage)
	auth.POST("/create", v1.CreateMenu)
	auth.PUT("/update/:id", v1.UpdateMenu)
	auth.DELETE("/:id", v1.DeleteMenu)
	auth.GET("/detail", v1.GetMenuDetail)
	auth.PUT("/changeState", v1.ChangeMenuState)
	auth.POST("/batchDelete", v1.BatchDeleteMenu)
}
