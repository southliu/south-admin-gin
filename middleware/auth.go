package middleware

import (
	"net/http"
	"south-admin-gin/models/system"
	"south-admin-gin/services/system"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.CommonResponse{
				Code: 401,
				Message: "未登录或 token 无效",
			})
			c.Abort()
			return
		}

		// 支持 "Bearer <token>" 和纯 token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := services.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.CommonResponse{
				Code: 401,
				Message: "token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息写入上下文，供后续接口使用
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
