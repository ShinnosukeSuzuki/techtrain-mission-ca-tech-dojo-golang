package controllers

type (
	// /user/create
	UserCreateRequest struct {
		Name string `json:"name"`
	}

	// /user/update
	UserUpdateRequest struct {
		Name string `json:"name"`
	}

	// /gacha/draw
	GachaDrawRequest struct {
		Times int `json:"times"`
	}
)
