package routers

import (
	"forum/handlers"
	"forum/utils"
)

func InitPagesRouter(app *utils.App) {
	pagesRouter := app.NewRouter("")

	pagesRouter.Get("/", handlers.HomeHandler)
	pagesRouter.Get("/register", handlers.RegisterHandler)
	pagesRouter.Get("/login", handlers.LoginHandler)
}
