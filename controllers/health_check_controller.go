package controllers

import (
	"encoding/json"

	"net/http"
)

// HealthCheck用のコントローラ構造体
type HealthCheckController struct{}

// コンストラクタ関数
func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

// ハンドラーメソッドを定義
// GET /health_check
func (c *HealthCheckController) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// GET以外のリクエストはエラー
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res := &HealthCheckResponse{
		Message: "OK",
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
