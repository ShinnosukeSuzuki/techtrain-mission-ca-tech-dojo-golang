package api

import (
	"database/sql"
	"net/http"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api/middleware"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/services"
)

func NewRouter(db *sql.DB) *http.ServeMux {

	// userRepositoryのインスタンスを生成
	uRep := repositories.NewUserRepository(db)

	// userServiceのインスタンスを生成
	uSer := services.NewUserService(uRep)

	// userControllerのインスタンスを生成
	uCon := controllers.NewUserController(uSer)

	// register routes
	mux := http.NewServeMux()

	// ルーティングの設定
	mux.Handle("/user/create", middleware.JSONContentTypeMiddleware(http.HandlerFunc(uCon.UserCreateHandler)))
	// /user/getと/user/updateではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/user/get", middleware.XTokenAuthMiddleware(middleware.JSONContentTypeMiddleware(http.HandlerFunc(uCon.UserGetHandler)), uRep))
	mux.Handle("/user/update", middleware.XTokenAuthMiddleware(http.HandlerFunc(uCon.UserUpdateHandler), uRep))

	return mux
}
