package routers

import (
	"net/http"

	"forum/handlers"
	"forum/handlers/api"
	"forum/handlers/middleware"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /assets/", handlers.AssetsHandler)

	// pages
	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("GET /login", middleware.RedirectMiddleware(handlers.LoginHandler))
	mux.HandleFunc("GET /register", middleware.RedirectMiddleware(handlers.RegisterHandler))
	mux.HandleFunc("GET /createpost", middleware.Middleware(handlers.CreatePostHandler))

	// api
	mux.HandleFunc("GET /api/posts", handlers.LoadData)
	mux.HandleFunc("GET /api/posts/{id}", handlers.LoadPostData)
	mux.HandleFunc("GET /api/categories", handlers.LoadAllCategories)

	// auth
	mux.HandleFunc("POST /api/login", api.LoginUser)
	mux.HandleFunc("POST /api/users", api.PostUser)

	// user activity
	mux.HandleFunc("POST /newposts", middleware.Middleware(handlers.NewPostHandler))
	mux.HandleFunc("POST /newcomment", middleware.Middleware(handlers.NewCommentHandler))
	mux.HandleFunc("POST /reactions", middleware.Middleware(handlers.LikeDislikeHandler))
	mux.HandleFunc("POST /newcategories", handlers.CreateCategoriesHandler)

	// logout
	mux.HandleFunc("POST /logout", api.LogoutUser)


	//getsession
	// mux.HandleFunc("GET /", api.SessionHandler)
}
