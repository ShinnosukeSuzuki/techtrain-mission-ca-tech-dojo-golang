package api

import (
	"database/sql"
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api/middleware"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/services"
)

func NewRouter(db *sql.DB) *http.ServeMux {

	// sevicesのインスタンスを生成
	ser := services.NewMyAppService(db)

	// コントローラのインスタンスを生成
	uc := controllers.NewUserController(ser)

	// register routes
	mux := http.NewServeMux()

	// ルーティングの設定
	mux.Handle("/user/create", http.HandlerFunc(uc.UserCreateHandler))
	// /user/getと/user/updateではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/user/get", middleware.XTokenAuthMiddleware(http.HandlerFunc(uc.UserGetHandler), db))
	mux.Handle("/user/update", middleware.XTokenAuthMiddleware(http.HandlerFunc(uc.UserUpdateHandler), db))

	return mux
}
