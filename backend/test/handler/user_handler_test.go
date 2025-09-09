package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- モックの定義 ---

// MockUserService は service.UserService のインターフェースを満たすモックオブジェクトです。
// これにより、実際のビジネスロジック（DBアクセスやJWT生成）をシミュレートします。
type MockUserService struct {
	mock.Mock
}

// Login はモックされたログインメソッドです。
// テストケースごとに「この引数で呼ばれたら、この値を返す」という振る舞いを定義できます。
func (m *MockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

// --- テスト本体 ---

func TestUserHandler_Login(t *testing.T) {
	// Ginをテストモードに設定
	gin.SetMode(gin.TestMode)

	// テストケースを定義
	testCases := []struct {
		name            string            // テストケース名
		requestBody     map[string]string // リクエストボディ
		mockUsername    string            // モックに渡すユーザー名
		mockPassword    string            // モックに渡すパスワード
		mockReturnToken string            // モックが返すトークン
		mockReturnError error             // モックが返すエラー
		expectedStatus  int               // 期待するHTTPステータスコード
		expectedBody    string            // 期待するレスポンスボディ
	}{
		{
			name:            "正常系: ログイン成功",
			requestBody:     map[string]string{"username": "testuser", "password": "password123"},
			mockUsername:    "testuser",
			mockPassword:    "password123",
			mockReturnToken: "mocked.jwt.token",
			mockReturnError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"token":"mocked.jwt.token"}`,
		},
		{
			name:            "異常系: 認証情報が無効",
			requestBody:     map[string]string{"username": "wronguser", "password": "wrongpassword"},
			mockUsername:    "wronguser",
			mockPassword:    "wrongpassword",
			mockReturnToken: "",
			mockReturnError: errors.New("invalid credentials"),
			expectedStatus:  http.StatusUnauthorized,
			expectedBody:    `{"error":"Invalid username or password"}`,
		},
		{
			name:           "異常系: リクエストボディが不正 (パスワード欠け)",
			requestBody:    map[string]string{"username": "testuser"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
		{
			name:           "異常系: リクエストボディが空",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- セットアップ ---
			mockService := new(MockUserService)

			// モックの振る舞いを設定（リクエストボディが正常な場合のみ）
			if tc.expectedStatus != http.StatusBadRequest {
				mockService.On("Login", tc.mockUsername, tc.mockPassword).Return(tc.mockReturnToken, tc.mockReturnError)
			}

			// ハンドラとルーターを初期化
			userHandler := handler.NewUserHandler(mockService)
			router := gin.Default()
			router.POST("/login", userHandler.Login)

			// リクエストボディをJSONに変換
			body, _ := json.Marshal(tc.requestBody)

			// HTTPリクエストをシミュレート
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// レスポンスを記録するためのレコーダーを作成
			w := httptest.NewRecorder()

			// --- 実行 ---
			router.ServeHTTP(w, req)

			// --- 検証 ---
			// ステータスコードを検証
			assert.Equal(t, tc.expectedStatus, w.Code)
			// レスポンスボディを検証
			assert.JSONEq(t, tc.expectedBody, w.Body.String())

			// モックが期待通りに呼ばれたか検証
			if tc.expectedStatus != http.StatusBadRequest {
				mockService.AssertExpectations(t)
			}
		})
	}
}
