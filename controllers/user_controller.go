package controllers

import (
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers/services"
	"github.com/labstack/echo/v4"
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
func (c *UserController) CreateHandler(ctx echo.Context) error {
	req := &UserCreateRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	user, err := c.uSer.Create(req.Name)
	if err != nil {
		return err
	}

	res := &UserCreateResponse{
		Token: user.Token,
	}

	return ctx.JSON(http.StatusOK, res)
}

// GET /user/get
func (c *UserController) GetHandler(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, "Invalid request")
	}

	user, err := c.uSer.Get(userID)
	if err != nil {
		return err
	}

	res := &UserGetResponse{
		Name: user.Name,
	}

	return ctx.JSON(http.StatusOK, res)
}

// PUT /user/update
func (c *UserController) UpdateNameHandler(ctx echo.Context) error {
	req := &UserUpdateRequest{}
	if err := ctx.Bind(req); err != nil {
		return err
	}

	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, "Invalid request")
	}

	if err := c.uSer.UpdateName(userID, req.Name); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
