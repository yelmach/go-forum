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

	// auth
	mux.HandleFunc("POST /api/login", api.LoginUser)
	mux.HandleFunc("POST /api/users", api.PostUser)

	// user activity
	mux.HandleFunc("POST /newposts", api.CreatePostsHandler)
	mux.HandleFunc("POST /newcomment", api.CreateCommentsHandler)
	mux.HandleFunc("POST /reactions", api.AddLikeDislikeHandler)
	mux.HandleFunc("POST /newcategories",api.CreateCategoriesHandler)

	// logout
	mux.HandleFunc("GET /logout", api.LogoutUser)
}
