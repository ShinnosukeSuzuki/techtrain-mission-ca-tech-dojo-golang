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

	// user関連
	uRep := repositories.NewUserRepository(db)
	uSer := services.NewUserService(uRep)
	uCon := controllers.NewUserController(uSer)

	// user_character関連
	ucRep := repositories.NewUserCharacterRepository(db)
	ucSer := services.NewUserCharacterService(ucRep)
	ucCon := controllers.NewUserCharacterController(ucSer)

	// register routes
	mux := http.NewServeMux()

	// ルーティングの設定
	mux.Handle("/user/create", middleware.JSONContentTypeMiddleware(http.HandlerFunc(uCon.UserCreateHandler)))
	// /user/getと/user/updateではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/user/get", middleware.XTokenAuthMiddleware(middleware.JSONContentTypeMiddleware(http.HandlerFunc(uCon.UserGetHandler)), uRep))
	mux.Handle("/user/update", middleware.XTokenAuthMiddleware(http.HandlerFunc(uCon.UserUpdateHandler), uRep))

	// /character/listではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/character/list", middleware.XTokenAuthMiddleware(middleware.JSONContentTypeMiddleware(http.HandlerFunc(ucCon.UserCharacterGetHandler)), uRep))

	return mux
}
