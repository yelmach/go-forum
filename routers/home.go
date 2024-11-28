package routers

import (
	"net/http"

	"forum/handlers"
	"forum/handlers/middleware"
)

// Setup Home Routes
func SetupHomeRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/login", middleware.RedirectMiddleware(handlers.LoginHandler))
	mux.HandleFunc("/register", middleware.RedirectMiddleware(handlers.RegisterHandler))

	mux.HandleFunc("/newpost", middleware.Middleware(handlers.NewPostHandler))
	mux.HandleFunc("/newcomment", middleware.Middleware(handlers.NewCommentHandler))
	mux.HandleFunc("/reaction", middleware.Middleware(handlers.ReactionHandler))
	return mux
}
