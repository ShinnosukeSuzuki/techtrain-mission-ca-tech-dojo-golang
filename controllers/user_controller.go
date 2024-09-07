package controllers

import (
	"encoding/json"

	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api/middleware"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
)

// User用のコントローラ構造体
type UserController struct {
	uSer services.UserServicer
}

// コンストラクタ関数
func NewUserController(s services.UserServicer) *UserController {
	return &UserController{uSer: s}
}

// ハンドラーメソッドを定義
// POST /user/create
func (c *UserController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// POST以外のリクエストはエラー
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := &UserCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.uSer.Create(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &UserCreateResponse{
		Token: user.Token,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GET /user/get
func (c *UserController) GetHandler(w http.ResponseWriter, r *http.Request) {
	// GET以外のリクエストはエラー
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := r.Context().Value(middleware.UserIDKeyType{}).(string)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := c.uSer.Get(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &UserGetResponse{
		Name: user.Name,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /user/update
func (c *UserController) UpdateNameHandler(w http.ResponseWriter, r *http.Request) {
	// PUT以外のリクエストはエラー
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := &UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(middleware.UserIDKeyType{}).(string)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := c.uSer.UpdateName(userId, req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
