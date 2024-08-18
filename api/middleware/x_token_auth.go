package middleware

import (
	"database/sql"
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/common"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

func XTokenAuthMiddleware(h http.Handler, db *sql.DB) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーからX-Tokenを取得
		xToken := r.Header.Get("X-Token")
		if xToken == "" {
			http.Error(w, "X-Token is required", http.StatusUnauthorized)
			return
		}

		// 取得したX-Tokenを持つユーザーが存在するか確認
		user, err := repositories.GetUserByToken(db, xToken)
		if err != nil || user.Token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// X-Tokenをcontextに保存
		r = common.SetToken(r, xToken)
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}
