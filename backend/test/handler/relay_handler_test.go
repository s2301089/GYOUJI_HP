package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
)

// テスト用モックRelayService

type mockRelayService struct{}

func (m *mockRelayService) GetRelayResults(relayType string) ([]model.RelayResult, error) {
	return []model.RelayResult{
		{ID: 1, RelayType: relayType, RelayRank: 1, ClassID: 1, CreatedAt: "2025-09-14"},
	}, nil
}

func (m *mockRelayService) RegisterRelayResults(relayType string, classIDs []int) error {
	return nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	// DB依存を排除し、モックRelayServiceをセット
	handler.RelayService = &mockRelayService{}
	r := gin.Default()
	r.GET("/api/relay", handler.GetRelayScores)
	r.POST("/api/relay", handler.PostRelayScores)
	return r
}

func TestGetRelayScores_Unauthorized(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/relay?relay_type=A", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
		t.Errorf("expected unauthorized/forbidden, got %d", w.Code)
	}
}

func TestPostRelayScores_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/relay", strings.NewReader("{")) // 不正なJSON
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user", map[string]interface{}{"role": "admin", "username": "admin_relay", "assigned_sport": "relay"})
	handler.PostRelayScores(c)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected bad request, got %d", w.Code)
	}
}

func TestGetRelayScores_Success(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/relay?relay_type=A", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user", map[string]interface{}{"role": "superroot", "username": "superroot", "assigned_sport": "relay"})
	handler.GetRelayScores(c)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var results []model.RelayResult
	_ = json.Unmarshal(w.Body.Bytes(), &results)
}

func TestPostRelayScores_Success(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"relay_type":"A","class_ids":[1,2,3,4,5,6]}`
	req, _ := http.NewRequest("POST", "/api/relay", strings.NewReader(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user", map[string]interface{}{"role": "admin", "username": "admin_relay", "assigned_sport": "relay"})
	handler.PostRelayScores(c)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
}
