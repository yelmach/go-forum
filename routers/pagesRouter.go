package routers

import (
	"forum/handlers"
	"forum/utils"
)

func InitPagesRouter(app *utils.App) {
	pagesRouter := app.NewRouter("")

	// Severing Assets (Don't worry that way of serving static file just for demo)
	pagesRouter.Get("/assets/css/main.css", handlers.AssetsHandler)
	pagesRouter.Get("/assets/css/register_login.css", handlers.AssetsHandler)
	pagesRouter.Get("/assets/img/logo.svg", handlers.AssetsHandler)
	pagesRouter.Get("/assets/img/favicon.ico", handlers.AssetsHandler)
	pagesRouter.Get("/assets/img/Image01.png", handlers.AssetsHandler)
	pagesRouter.Get("/assets/img/Image02.png", handlers.AssetsHandler)
	pagesRouter.Get("/assets/img/usrs/Ava01.png", handlers.AssetsHandler)

	pagesRouter.Get("/", handlers.HomeHandler)
	pagesRouter.Get("/register", handlers.RegisterHandler)
	pagesRouter.Get("/login", handlers.LoginHandler)
}