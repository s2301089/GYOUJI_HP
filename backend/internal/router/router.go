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
func SetupRouter(userHandler *handler.UserHandler, tournamentHandler *handler.TournamentHandler, matchHandler *handler.MatchHandler, jwtSecret string, scoreHandler *handler.ScoreHandler) *gin.Engine {
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
		// スコア取得
		api.GET("/score", scoreHandler.GetScores)

		// --- 認証関連のルートグループ ---
		auth := api.Group("/auth")
		{
			// ログインは認証不要
			auth.POST("/login", userHandler.Login)

			// ログアウト・ユーザー情報取得は認証必須
			auth.Use(middleware.AuthMiddleware(jwtSecret))
			{
				auth.POST("/logout", userHandler.Logout)
				auth.GET("/me", userHandler.GetMe)
			}
		}

		// 試合スコア更新API（認証必須）
		api.PUT("/matches/:id", middleware.AuthMiddleware(jwtSecret), matchHandler.UpdateMatchScore)
		// 試合一覧取得API（認証必須）
		api.GET("/matches/:sport", middleware.AuthMiddleware(jwtSecret), matchHandler.GetMatchesBySport)
		// リレー関連
		api.GET("/relay", middleware.AuthMiddleware(jwtSecret), handler.GetRelayScores)
		api.POST("/relay", middleware.AuthMiddleware(jwtSecret), handler.PostRelayScores)
	}

	return r
}
