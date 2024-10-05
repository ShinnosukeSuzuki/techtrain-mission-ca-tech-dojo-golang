package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
)

// GachaDraw用のコントローラ構造体
type GachaDrawController struct {
	service services.GachaDrawServicer
}

// コンストラクタ関数
func NewGachaDrawController(s services.GachaDrawServicer) *GachaDrawController {
	return &GachaDrawController{service: s}
}

// ハンドラーメソッドを定義
// POST /gacha/draw
func (c *GachaDrawController) DrawHandler(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, "Invalid request")
	}

	// リクエストボディをパース
	var req GachaDrawRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// timesが0未満の場合はエラー
	if req.Times < 0 {
		return ctx.JSON(http.StatusBadRequest, "Invalid request")
	}

	// ガチャを引く
	results, err := c.service.Draw(req.Times, userID)
	if err != nil {
		return err
	}

	gachaResults := make([]GachaResult, 0, len(results))
	for _, character := range results {
		gachaResults = append(gachaResults, GachaResult{
			CharacterID: character.ID,
			Name:        character.Name,
		})
	}

	res := &GachaDrawResponse{
		Results: gachaResults,
	}

	return ctx.JSON(http.StatusOK, res)
}
