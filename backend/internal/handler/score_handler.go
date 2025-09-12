package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type ScoreHandler struct{ Service service.ScoreServiceInterface }

func NewScoreHandler(s service.ScoreServiceInterface) *ScoreHandler { return &ScoreHandler{Service: s} }

func (h *ScoreHandler) GetScores(c *gin.Context) {
	scores, err := h.Service.GetScores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scores)
}
