package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api/middleware"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// UserCharacter用のコントローラ構造体
type UserCharacterController struct {
	service services.UserCharacterServicer
}

// コンストラクタ関数
func NewUserCharacterController(s services.UserCharacterServicer) *UserCharacterController {
	return &UserCharacterController{service: s}
}

// ハンドラーメソッドを定義

// GET /character/list
func (c *UserCharacterController) GetListHandler(w http.ResponseWriter, r *http.Request) {
	// GET以外のリクエストはエラー
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, _ := r.Context().Value(middleware.UserIDKeyType{}).(string)

	// userIdを元に一致するユーザーが所持するキャラクターを取得
	characterList, err := c.service.List(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(characterList.Characters) == 0 {
		characterList.Characters = []models.UserCharacter{}
	}

	res := &CharacterListResponse{
		Characters: characterList.Characters,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
