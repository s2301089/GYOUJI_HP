package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuperrootOnly は superroot ロールを持つユーザーのみを許可するミドルウェアです。
func SuperrootOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "superroot" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: You don't have access to this resource"})
			return
		}
		c.Next()
	}
}
