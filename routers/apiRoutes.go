package routers

import (
	"forum/handlers/api"
	"forum/utils"
)

func InitApiRouter(app *utils.App) {
	apiRouter := app.NewRouter("/api")

	apiRouter.Post("/users", api.PostUser)
	apiRouter.Post("/login", api.LoginUser)
	apiRouter.Get("/check_session", api.SessionHandler)
}
