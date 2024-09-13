package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// contextに格納するUserIDのキー
type UserIDKeyType struct{}

// ログに出力する構造体を定義
type AccessLogging struct {
	Timestamp string `json:"timestamp"`
	UserID    string `json:"user_id,omitempty"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

func XTokenAuthMiddleware(h http.Handler, uRep repositories.UserRepository) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーからX-Tokenを取得
		xToken := r.Header.Get("X-Token")
		if xToken == "" {
			logAccess(r, "", http.StatusUnauthorized, "X-Token is required")
			http.Error(w, "X-Token is required", http.StatusUnauthorized)
			return
		}

		// 取得したX-Tokenを持つユーザーが存在するか確認
		user, err := uRep.GetByToken(xToken)
		if err != nil || user.Token == "" {
			logAccess(r, "", http.StatusUnauthorized, "Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logAccess(r, user.ID, http.StatusOK, "Authorized")

		// contextにUserIDを格納
		ctx := context.WithValue(r.Context(), UserIDKeyType{}, user.ID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}

func logAccess(r *http.Request, userID string, status int, message string) {
	al := AccessLogging{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000000 -0700 MST"),
		UserID:    userID,
		Path:      r.URL.Path,
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
