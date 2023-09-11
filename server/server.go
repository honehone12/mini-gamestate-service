package server

import (
	"mini-gamestate-service/db"
	"mini-gamestate-service/server/context"
	"mini-gamestate-service/server/handlers"
	"mini-gamestate-service/server/handlers/jewels"
	"mini-gamestate-service/server/handlers/sessions"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Run(
	name string,
	version string,
	listenAt string,
	db db.Orm,
) {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &context.Context{
				Context: c,
				Orm:     db,
				Metadata: context.Metadata{
					Name:    name,
					Version: version,
				},
			}
			return next(ctx)
		}
	})
	e.Validator = context.NewValidator()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/", handlers.Root)
	e.POST("/session/set", sessions.Set)
	e.POST("/session/verify", sessions.Verify)
	e.POST("/jewel/init", jewels.Initialize)
	e.POST("/jewel/incr", jewels.Incr)
	e.POST("/jewel/get-all", jewels.GetAll)

	e.Logger.SetLevel(log.WARN)
	e.Logger.Fatal(e.Start(listenAt))
}
