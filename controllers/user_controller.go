package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/common"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
)

// User用のコントローラ構造体
type UserController struct {
	service services.UserServicer
}

// コンストラクタ関数
func NewUserController(s services.UserServicer) *UserController {
	return &UserController{service: s}
}

// ユーザーのリクエスト・レスポンスボディを定義
type (
	// /user/createのリクエストボディ
	UserCreateRequest struct {
		Name string `json:"name"`
	}
	// /user/createのレスポンスボディ
	UserCreateResponse struct {
		Token string `json:"token"`
	}

	// /user/getのレスポンスボディ
	UserGetResponse struct {
		Name string `json:"name"`
	}

	// /user/updateのリクエストボディ
	UserUpdateRequest struct {
		Name string `json:"name"`
	}
)

// ハンドラーメソッドを定義
// POST /user/create
func (c *UserController) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	// POST以外のリクエストはエラー
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// リクエストボディをパース
	req := &UserCreateRequest{}
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
	res := &UserCreateResponse{
		Token: user.Token,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GET /user/get
func (c *UserController) UserGetHandler(w http.ResponseWriter, r *http.Request) {
	// GET以外のリクエストはエラー
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// X-Tokenをcontextから取得
	xToken := common.GetToken(r)

	// X-Tokenを持つユーザーを取得
	user, err := c.service.UserGetService(xToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを返却
	res := &UserGetResponse{
		Name: user.Name,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PUT /user/update
func (c *UserController) UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// PUT以外のリクエストはエラー
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// リクエストボディをパース
	req := &UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// X-Tokenをcontextから取得
	xToken := common.GetToken(r)

	// ユーザーのnameを更新
	if err := c.service.UserUpdateService(xToken, req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを返却(200を返すだけ)
	w.WriteHeader(http.StatusOK)
}
