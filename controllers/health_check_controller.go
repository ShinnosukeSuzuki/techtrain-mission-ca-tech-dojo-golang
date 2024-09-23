package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck用のコントローラ構造体
type HealthCheckController struct{}

// コンストラクタ関数
func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

// ハンドラーメソッドを定義
// GET /health-check
func (c *HealthCheckController) HealthCheckHandler(ctx echo.Context) error {

	res := &HealthCheckResponse{
		Message: "OK",
	}
	if err := ctx.Bind(res); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}
