package middleware

import (
	"context"
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

// contextに格納するUserIDのキー
type UserIDKeyType struct{}

func XTokenAuthMiddleware(h http.Handler, uRep repositories.UserRepository) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーからX-Tokenを取得
		xToken := r.Header.Get("X-Token")
		if xToken == "" {
			http.Error(w, "X-Token is required", http.StatusUnauthorized)
			return
		}

		// 取得したX-Tokenを持つユーザーが存在するか確認
		user, err := uRep.GetByToken(xToken)
		if err != nil || user.Token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// contextにUserIDを格納
		ctx := context.WithValue(r.Context(), UserIDKeyType{}, user.ID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}
