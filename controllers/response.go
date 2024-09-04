package controllers

import "github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"

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
	CharacterListResponse struct {
		Characters []models.UserCharacter `json:"characters"`
	}
)
