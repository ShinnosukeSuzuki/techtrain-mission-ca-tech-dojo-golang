package controllers

// ユーザーのリクエストボディを定義
type (
	// /user/create
	UserCreateRequest struct {
		Name string `json:"name"`
	}

	// /user/update
	UserUpdateRequest struct {
		Name string `json:"name"`
	}
)
