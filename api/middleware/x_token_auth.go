package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// ログに出力する構造体を定義
type AccessLogging struct {
	Timestamp string `json:"timestamp"`
	UserID    string `json:"user_id,omitempty"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

func XTokenAuthMiddleware(uRep repositories.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// リクエストヘッダーからX-Tokenを取得
			xToken := ctx.Request().Header.Get("X-Token")
			if xToken == "" {
				logAccess(ctx, "", http.StatusUnauthorized, "X-Token is required")
				return ctx.JSON(http.StatusUnauthorized, "X-Token is required")
			}

			// 取得したX-Tokenを持つユーザーが存在するか確認
			user, err := uRep.GetByToken(xToken)
			if err != nil || user.Token == "" {
				logAccess(ctx, "", http.StatusUnauthorized, "Unauthorized")
				return ctx.JSON(http.StatusUnauthorized, "Unauthorized")
			}

			logAccess(ctx, user.ID, http.StatusOK, "Authorized")

			// コンテキストにUserIDを格納
			ctx.Set("userID", user.ID)

			return next(ctx)
		}
	}
}

func logAccess(ctx echo.Context, userID string, status int, message string) {
	al := AccessLogging{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000000 -0700 MST"),
		UserID:    userID,
		Path:      ctx.Request().URL.Path,
		Status:    status,
		Message:   message,
	}

	// JSONに変換して出力
	jsonLog, err := json.Marshal(al)
	if err != nil {
		fmt.Println("Error marshaling log:", err)
		return
	}
	fmt.Println(string(jsonLog))
}
