package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

type TournamentHandler struct {
	Service service.TournamentServiceInterface
}

func NewTournamentHandler(s service.TournamentServiceInterface) *TournamentHandler {
	return &TournamentHandler{Service: s}
}

// GET /api/tournaments/:sport
func (h *TournamentHandler) GetTournamentsBySport(c *gin.Context) {
	sport := c.Param("sport")
	result, err := h.Service.GetTournamentsBySport(sport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
