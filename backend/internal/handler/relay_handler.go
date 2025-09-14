package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"
)

var RelayService service.RelayService

func InitRelayService(repo *repository.RelayRepository) {
	RelayService = service.NewRelayService(repo)
}

// GetRelayScores 学年対抗リレーの得点・結果取得API
// @Summary 学年対抗リレーの得点・順位取得
// @Tags relay
// @Param relay_type query string true "リレー種別 (A or B)"
// @Success 200 {array} model.RelayResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/relay [get]
func GetRelayScores(c *gin.Context) {
	// 認証ユーザー情報取得
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, _ := userMap["role"].(string)
	username, _ := userMap["username"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)
	// superrootまたはadmin_relayのみ許可
	if !(role == "superroot" || (role == "admin" && username == "admin_relay" && assignedSport == "relay")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	relayType := c.Query("relay_type") // 'A' or 'B'
	if relayType != "A" && relayType != "B" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "relay_type must be 'A' or 'B'"})
		return
	}
	results, err := RelayService.GetRelayResults(relayType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// 学年対抗リレーの得点・結果登録API
// PostRelayScores 学年対抗リレーの得点・結果登録API
// @Summary 学年対抗リレーの得点・順位登録
// @Tags relay
// @Param body body RelayRegisterRequest true "順位順のclass_id配列とリレー種別"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/relay [post]
func PostRelayScores(c *gin.Context) {
	// 認証ユーザー情報取得
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userMap, ok := user.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, _ := userMap["role"].(string)
	username, _ := userMap["username"].(string)
	assignedSport, _ := userMap["assigned_sport"].(string)
	// superrootまたはadmin_relayのみ許可
	if !(role == "superroot" || (role == "admin" && username == "admin_relay" && assignedSport == "relay")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req model.RelayRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.RelayType != "A" && req.RelayType != "B" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "relay_type must be 'A' or 'B'"})
		return
	}
	if len(req.ClassIDs) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_ids must be 6 items (順位順)"})
		return
	}
	err := RelayService.RegisterRelayResults(req.RelayType, req.ClassIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "登録完了"})
}
