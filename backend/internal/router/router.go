package router

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter は Gin のルーターをセットアップし、ルートを定義します。
func SetupRouter(userHandler *handler.UserHandler, jwtSecret string) *gin.Engine {
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
		}

		// --- 認証が必要なルートグループ ---
		authRequired := api.Group("/")
		// このグループのルートは、すべてAuthMiddlewareを通過する必要がある
		authRequired.Use(middleware.AuthMiddleware(jwtSecret))
		{
			authRequired.POST("/auth/logout", userHandler.Logout)
			// 今後、ユーザー情報の取得など、認証が必要なAPIはここに追加します
		}
	}

	return r
}
