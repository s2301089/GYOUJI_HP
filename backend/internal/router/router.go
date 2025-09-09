package router

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter は Gin のルーターをセットアップし、ルートを定義します。
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// CORSミドルウェアの設定（開発用に寛容な設定）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			// ここに /logout などのエンドポイントを追加可能
		}
	}

	return r
}
