package api

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api/middleware"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/cache"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/controllers"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/repositories"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/services"
	"github.com/labstack/echo/v4"
)

func NewRouter(db *sql.DB, characterCache *cache.CharacterProbabilityCache) *echo.Echo {
	e := echo.New()

	// Repositories
	uRep := repositories.NewUserRepository(db)
	ucRep := repositories.NewUserCharacterRepository(db)

	// Services
	uSer := services.NewUserService(uRep)
	ucSer := services.NewUserCharacterService(ucRep, characterCache)
	gdSer := services.NewGachaDrawService(ucRep, characterCache)

	// Controllers
	hCon := controllers.NewHealthCheckController()
	uCon := controllers.NewUserController(uSer)
	ucCon := controllers.NewUserCharacterController(ucSer)
	gdCon := controllers.NewGachaDrawController(gdSer)

	// Routes
	e.GET("/health-check", hCon.HealthCheckHandler)
	e.POST("/user/create", uCon.CreateHandler)

	// Routes that require X-Token authentication
	authGroup := e.Group("", middleware.XTokenAuthMiddleware(uRep))
	authGroup.GET("/user/get", uCon.GetHandler)
	authGroup.PUT("/user/update", uCon.UpdateNameHandler)
	authGroup.POST("/gacha/draw", gdCon.DrawHandler)
	authGroup.GET("/character/list", ucCon.GetListHandler)

	return e
}
