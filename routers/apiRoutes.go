package routers

import (
	"net/http"

	"forum/handlers"
	"forum/handlers/api"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /assets/", handlers.AssetsHandler)

	// pages
	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /login", handlers.LoginHandler)
	mux.HandleFunc("GET /register", handlers.RegisterHandler)
	mux.HandleFunc("GET /createpost", handlers.CreatePostHandler)

	// api
	mux.HandleFunc("GET /api/posts", handlers.LoadData)
	mux.HandleFunc("GET /api/posts/{id}", handlers.LoadPostData)
	mux.HandleFunc("GET /api/categories", handlers.LoadAllCategories)

	// auth
	mux.HandleFunc("POST /api/login", api.LoginUser)
	mux.HandleFunc("POST /api/users", api.PostUser)

	// user activity
	mux.HandleFunc("POST /newposts", handlers.NewPostHandler)
	mux.HandleFunc("POST /", handlers.NewCommentHandler)
	mux.HandleFunc("POST /reactions", handlers.LikeDislikeHandler)
	mux.HandleFunc("POST /newcategories", handlers.CreateCategoriesHandler)

	// logout
	mux.HandleFunc("POST /logout", api.LogoutUser)
}
