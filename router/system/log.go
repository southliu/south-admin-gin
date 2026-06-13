package system

import (
	v1 "serve-wechat-gin/api/v1/system"

	"github.com/gin-gonic/gin"
)

func InitLogRoutes(r *gin.RouterGroup) {
	r.POST("/create", v1.CreateLog)

	auth := withAuth(r)
	auth.GET("/page", v1.GetLogPage)
	auth.DELETE("/:id", v1.DeleteLog)
	auth.POST("/batchDelete", v1.BatchDeleteLog)
}
