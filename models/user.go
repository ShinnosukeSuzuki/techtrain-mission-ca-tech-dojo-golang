package models

type (
	// ユーザーの構造体
	User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}

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
)
