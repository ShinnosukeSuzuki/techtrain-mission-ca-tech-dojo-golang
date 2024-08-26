package controllers

// ユーザーのレスポンスボディを定義
type (
	// /user/create
	UserCreateResponse struct {
		Token string `json:"token"`
	}

	// /user/get
	UserGetResponse struct {
		Name string `json:"name"`
	}
)
