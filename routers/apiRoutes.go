package routers

import (
	"net/http"

	"forum/handlers"
	"forum/handlers/api"
	"forum/handlers/auth"
	"forum/handlers/middleware"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/assets/", handlers.AssetsHandler)

	// pages
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/login", middleware.RedirectMiddleware(handlers.LoginHandler))
	mux.HandleFunc("/register", middleware.RedirectMiddleware(handlers.RegisterHandler))

	// api
	mux.HandleFunc("/api/posts", api.LoadData)
	mux.HandleFunc("/api/posts/{id}", api.LoadPostData)
	mux.HandleFunc("/api/categories", api.LoadAllCategories)

	// auth
	mux.HandleFunc("/auth/register", auth.RegisterUser)
	mux.HandleFunc("/auth/login", auth.LoginUser)
	mux.HandleFunc("/auth/logout", auth.LogOutUser)

	// user activity
	mux.HandleFunc("/newpost", middleware.Middleware(handlers.NewPostHandler))
	mux.HandleFunc("/newcomment", middleware.Middleware(handlers.NewCommentHandler))
	mux.HandleFunc("/reaction", middleware.Middleware(handlers.ReactionHandler))
}
