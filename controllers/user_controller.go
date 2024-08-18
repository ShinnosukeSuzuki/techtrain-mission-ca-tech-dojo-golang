package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
)

// User用のコントローラ構造体
type UserController struct {
	service services.UserServicer
}

// コンストラクタ関数
func NewUserController(s services.UserServicer) *UserController {
	return &UserController{service: s}
}

// ハンドラーメソッドを定義
// POST /user/create
func (c *UserController) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	// POST以外のリクエストはエラー
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// リクエストボディをパース
	req := &models.UserCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ユーザーを作成
	user, err := c.service.UserCreateService(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを返却
	res := &models.UserCreateResponse{
		Token: user.Token,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
