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

	// health_check関連
	hCon := controllers.NewHealthCheckController()

	// user関連
	uRep := repositories.NewUserRepository(db)
	uSer := services.NewUserService(uRep)
	uCon := controllers.NewUserController(uSer)

	// user_character関連
	ucRep := repositories.NewUserCharacterRepository(db)
	ucSer := services.NewUserCharacterService(ucRep)
	ucCon := controllers.NewUserCharacterController(ucSer)

	// gacha draw関連
	cRep := repositories.NewCharacterRepository(db)
	gdSer := services.NewGachaDrawService(ucRep, cRep)
	gdCon := controllers.NewGachaDrawController(gdSer)

	// register routes
	mux := http.NewServeMux()

	// ルーティングの設定
	mux.Handle("/health-check", middleware.JSONContentTypeMiddleware(http.HandlerFunc(hCon.HealthCheckHandler)))
	mux.Handle("/user/create", middleware.JSONContentTypeMiddleware(http.HandlerFunc(uCon.CreateHandler)))
	// /user/getと/user/updateではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/user/get", middleware.XTokenAuthMiddleware(middleware.JSONContentTypeMiddleware(http.HandlerFunc(uCon.GetHandler)), uRep))
	mux.Handle("/user/update", middleware.XTokenAuthMiddleware(http.HandlerFunc(uCon.UpdateNameHandler), uRep))

	// /gacha/drawではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/gacha/draw", middleware.XTokenAuthMiddleware(middleware.JSONContentTypeMiddleware(http.HandlerFunc(gdCon.DrawHandler)), uRep))

	// /character/listではX-Tokenが必要なのでmiddlewareを適用
	mux.Handle("/character/list", middleware.XTokenAuthMiddleware(middleware.JSONContentTypeMiddleware(http.HandlerFunc(ucCon.GetListHandler)), uRep))

	return mux
}
