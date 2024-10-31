package routers

import (
	"forum/handlers"
	"forum/handlers/api"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/assets/", handlers.AssetsHandler)

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	mux.HandleFunc("/api/login", api.LoginUser)
	mux.HandleFunc("/api/users", api.PostUser)
	mux.HandleFunc("/check_session", api.SessionHandler)
	// mux.HandleFunc("POST /post")
	// mux.HandleFunc("POST /comment")
	// mux.HandleFunc("POST /check_session")
}
