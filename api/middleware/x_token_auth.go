package middleware

import (
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/common"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
)

func XTokenAuthMiddleware(h http.Handler, rep repositories.UserRepository) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーからX-Tokenを取得
		xToken := r.Header.Get("X-Token")
		if xToken == "" {
			http.Error(w, "X-Token is required", http.StatusUnauthorized)
			return
		}

		// 取得したX-Tokenを持つユーザーが存在するか確認
		user, err := rep.GetUserByToken(xToken)
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
