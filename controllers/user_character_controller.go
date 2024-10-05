package controllers

import (
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/models"
	"github.com/labstack/echo/v4"
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
func (c *UserCharacterController) GetListHandler(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, "Invalid request")
	}

	// userIdを元に一致するユーザーが所持するキャラクターを取得
	characterList, err := c.service.List(userID)
	if err != nil {
		return err
	}

	if len(characterList.Characters) == 0 {
		characterList.Characters = []models.UserCharacter{}
	}

	characters := make([]UserCharacter, 0, len(characterList.Characters))
	for _, character := range characterList.Characters {
		characters = append(characters, UserCharacter{
			UserCharacterID: character.UserCharacterID,
			CharacterID:     character.CharacterID,
			Name:            character.Name,
		})
	}

	res := &CharacterListResponse{
		Characters: characters,
	}

	return ctx.JSON(http.StatusOK, res)
}
