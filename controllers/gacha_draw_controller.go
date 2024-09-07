package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api/middleware"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
)

// GachaDraw用のコントローラ構造体
type GachaDrawController struct {
	service services.GachaDrawServicer
}

// コンストラクタ関数
func NewGachaDrawController(s services.GachaDrawServicer) *GachaDrawController {
	return &GachaDrawController{service: s}
}

// ハンドラーメソッドを定義
// POST /gacha/draw
func (c *GachaDrawController) DrawHandler(w http.ResponseWriter, r *http.Request) {
	// POST以外のリクエストはエラー
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := r.Context().Value(middleware.UserIDKeyType{}).(string)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// リクエストボディをパース
	var req GachaDrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// timesが0未満の場合はエラー
	if req.Times < 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ガチャを引く
	results, err := c.service.Draw(req.Times, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gachaResults := make([]GachaResult, 0, len(results))
	for _, character := range results {
		gachaResults = append(gachaResults, GachaResult{
			CharacterID: character.CharacterID,
			Name:        character.Name,
		})
	}

	res := &GachaDrawResponse{
		Results: gachaResults,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
