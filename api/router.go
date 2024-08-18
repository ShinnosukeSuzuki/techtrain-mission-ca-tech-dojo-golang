package api

import (
	"database/sql"
	"net/http"

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
	mux.HandleFunc("/user/create", uc.UserCreateHandler)

	return mux
}
