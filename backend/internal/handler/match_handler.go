package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type MatchHandler struct {
	Service service.MatchServiceInterface
}

func NewMatchHandler(s service.MatchServiceInterface) *MatchHandler {
	return &MatchHandler{Service: s}
}

// PUT /api/matches/:id
// UpdateMatchScore godoc
// @Summary      試合のスコアを更新
// @Description  指定されたIDの試合のスコアを更新し、勝者を次の試合へ進出させます。superrootまたは担当競技のadmin権限が必要です。
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "試合ID"
// @Param        score body      model.TeamScore true  "更新するスコア情報"
// @Param        Authorization header string true "Bearerトークン"
// @Success      200  {object}  model.Match "更新後の試合情報"
// @Failure      400  {object}  model.ErrorResponse "無効なリクエストです (Invalid match ID or request body)"
// @Failure      401  {object}  model.ErrorResponse "認証が必要です (Unauthorized)"
// @Failure      403  {object}  model.ErrorResponse "権限がありません (Forbidden)"
// @Failure      404  {object}  model.ErrorResponse "試合が見つかりません (Match not found)"
// @Failure      500  {object}  model.ErrorResponse "サーバー内部エラー"
// @Router       /api/matches/{id} [put]
// @Security     ApiKeyAuth
func (h *MatchHandler) UpdateMatchScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	var req model.TeamScore
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, exists := c.Get("user") // AuthMiddlewareでセットされている前提
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result, err := h.Service.UpdateMatchScore(id, *req.Team1Score, *req.Team2Score, user)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this match."})
			return
		} else if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Match not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
