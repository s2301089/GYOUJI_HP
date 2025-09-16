package handler

import (
	"net/http"

	"github.com/saku0512/GYOUJI_HP/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingHandler は設定関連のHTTPリクエストを処理します。
type SettingHandler struct {
	settingService service.SettingService
}

// NewSettingHandler は SettingHandler の新しいインスタンスを生成します。
func NewSettingHandler(s service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: s}
}

// GetVisibility はスコア表示設定を取得します。
func (h *SettingHandler) GetVisibility(c *gin.Context) {
	value, err := h.settingService.GetVisibility()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get setting"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"showTotalScores": value})
}

// UpdateVisibility はスコア表示設定を更新します。
func (h *SettingHandler) UpdateVisibility(c *gin.Context) {
	var req struct {
		ShowTotalScores bool `json:"showTotalScores"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.settingService.UpdateVisibility(req.ShowTotalScores); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully"})
}
