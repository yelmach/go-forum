package routers

import (
	"net/http"

	"forum/handlers"
	"forum/handlers/api"
	"forum/handlers/auth"
	"forum/handlers/middleware"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /assets/", handlers.AssetsHandler)

	// pages
	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /login", middleware.RedirectMiddleware(handlers.LoginHandler))
	mux.HandleFunc("GET /register", middleware.RedirectMiddleware(handlers.RegisterHandler))

	// api
	mux.HandleFunc("GET /api/posts", api.LoadData)
	mux.HandleFunc("GET /api/posts/{id}", api.LoadPostData)
	mux.HandleFunc("GET /api/categories", api.LoadAllCategories)

	// auth
	mux.HandleFunc("POST /auth/register", auth.RegisterUser)
	mux.HandleFunc("POST /auth/login", auth.LoginUser)
	mux.HandleFunc("POST /auth/logout", auth.LogOutUser)

	// user activity
	mux.HandleFunc("POST /newpost", middleware.Middleware(handlers.NewPostHandler))
	mux.HandleFunc("POST /newcomment", middleware.Middleware(handlers.NewCommentHandler))
	mux.HandleFunc("POST /reaction", middleware.Middleware(handlers.ReactionHandler))
}
