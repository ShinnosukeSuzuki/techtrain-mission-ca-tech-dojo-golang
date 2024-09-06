package controllers

type (
	// /user/create
	UserCreateResponse struct {
		Token string `json:"token"`
	}

	// /user/get
	UserGetResponse struct {
		Name string `json:"name"`
	}

	// /character/list
	UserCharacter struct {
		UserCharacterID string `json:"userCharacterID"`
		CharacterID     string `json:"characterID"`
		Name            string `json:"name"`
	}
	CharacterListResponse struct {
		Characters []UserCharacter `json:"characters"`
	}
)
