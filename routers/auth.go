package routers

import (
	"net/http"

	"forum/handlers/auth"
)

// Setup Auth Routes
func SetupAuthRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", auth.RegisterUser)
	mux.HandleFunc("/auth/login", auth.LoginUser)
	mux.HandleFunc("/auth/logout", auth.LogoutUser)

	return mux
}
