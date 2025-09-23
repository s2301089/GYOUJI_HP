package router

import (
	"os"
	"strings"

	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/middleware"

	"github.com/gin-gonic/gin"

	_ "github.com/saku0512/GYOUJI_HP/backend/docs"
)

// securityHeadersMiddleware は、推奨されるセキュリティヘッダーをレスポンスに追加します。
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

// SetupRouter は Gin のルーターをセットアップし、ルートを定義します。
func SetupRouter(userHandler *handler.UserHandler, tournamentHandler *handler.TournamentHandler, matchHandler *handler.MatchHandler, jwtSecret string, scoreHandler *handler.ScoreHandler, relayHandler *handler.RelayHandler, attendanceHandler *handler.AttendanceHandler, settingHandler *handler.SettingHandler) *gin.Engine {
	r := gin.Default()

	// セキュリティヘッダーミドルウェアを適用
	r.Use(securityHeadersMiddleware())

	// CORSミドルウェアの設定
	r.Use(func(c *gin.Context) {
		allowedOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOriginsStr == "" {
			// 環境変数が設定されていない場合は、開発用のデフォルト値を使用
			allowedOriginsStr = "http://localhost:5173,http://localhost:3300"
		}
		allowedOrigins := strings.Split(allowedOriginsStr, ",")

		origin := c.Request.Header.Get("Origin")
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == strings.TrimSpace(allowedOrigin) {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-control-allow-methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		// --- 認証が不要なルート ---
		api.GET("/tournaments/:sport", tournamentHandler.GetTournamentsBySport)
		// スコア取得
		api.GET("/score", scoreHandler.GetScores)

		// --- 認証関連のルートグループ ---
		auth := api.Group("/auth")
		{
			// ログイン・ログアウトは認証不要
			auth.POST("/login", userHandler.Login)
			auth.POST("/logout", userHandler.Logout)

			// ユーザー情報取得は認証必須
			auth.Use(middleware.AuthMiddleware(jwtSecret))
			{
				auth.GET("/me", userHandler.GetMe)
			}
		}

		// Settings
		api.GET("/settings/visibility", middleware.AuthMiddleware(jwtSecret), settingHandler.GetVisibility)
		api.PUT("/settings/visibility", middleware.AuthMiddleware(jwtSecret), middleware.SuperrootOnly(), settingHandler.UpdateVisibility)
		api.GET("/settings/weather", middleware.AuthMiddleware(jwtSecret), settingHandler.GetWeather)
		api.PUT("/settings/weather", middleware.AuthMiddleware(jwtSecret), middleware.TableTennisAdminOnly(), settingHandler.UpdateWeather)

		// 試合スコア更新API（認証必須）
		api.PUT("/matches/:id", middleware.AuthMiddleware(jwtSecret), matchHandler.UpdateMatchScore)
		// 試合一覧取得API（認証必須）
		api.GET("/matches/:sport", middleware.AuthMiddleware(jwtSecret), matchHandler.GetMatchesBySport)
		// 試合リセットAPI (superrootのみ)
		api.POST("/matches/:id/reset", middleware.AuthMiddleware(jwtSecret), middleware.SuperrootOnly(), matchHandler.ResetMatch)
		// リレー関連
		api.GET("/relay", relayHandler.GetRelayRankings)
		api.POST("/relay", middleware.AuthMiddleware(jwtSecret), relayHandler.RegisterRelayRankings)
		api.POST("/relay/reset", middleware.AuthMiddleware(jwtSecret), relayHandler.ResetRelay)

		// 出席点関連（superrootのみ）
		api.GET("/attendance", middleware.AuthMiddleware(jwtSecret), attendanceHandler.GetAttendanceScores)
		api.PUT("/attendance/:class_id", middleware.AuthMiddleware(jwtSecret), attendanceHandler.UpdateAttendanceScore)
		api.PUT("/attendance/batch", middleware.AuthMiddleware(jwtSecret), attendanceHandler.BatchUpdateAttendanceScores)
	}

	return r
}
