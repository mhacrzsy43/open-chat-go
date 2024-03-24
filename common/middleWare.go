package common

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenExtractionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取Token
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// 解析Token（示例逻辑）
		token := strings.TrimPrefix(tokenHeader, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// 验证Token的逻辑可以在这里实现

		// 将解析后的Token添加到Gin的Context中，以供后续的处理器使用
		c.Set("token", token)

		c.Next() // 调用下一个中间件或处理器
	}
}
