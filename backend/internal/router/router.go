package router

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/saku0512/GYOUJI_HP/backend/docs"
)

// SetupRouter は Gin のルーターをセットアップし、ルートを定義します。
func SetupRouter(userHandler *handler.UserHandler, tournamentHandler *handler.TournamentHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// CORSミドルウェアの設定（開発用に寛容な設定）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-control-allow-methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		// --- 認証が不要なルート ---
		api.GET("/tournaments/:sport", tournamentHandler.GetTournamentsBySport)

		// --- 認証関連のルートグループ ---
		auth := api.Group("/auth")
		{
			// ログインは認証不要
			auth.POST("/login", userHandler.Login)

			// ログアウトは認証必須
			auth.Use(middleware.AuthMiddleware(jwtSecret))
			{
				auth.POST("/logout", userHandler.Logout)
			}
		}

		// 今後、ユーザー情報の取得など、認証が必要なAPIは以下のようなグループに追加します
		// protected := api.Group("/protected")
		// protected.Use(middleware.AuthMiddleware(jwtSecret))
		// {
		// 	// protected.GET("/users/me", ...)
		// }
	}

	return r
}
