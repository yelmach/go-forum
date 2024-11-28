package routers

import "net/http"

// SetupRoutes initialize all endpoints of forum project
func SetupRoutes(rootMux *http.ServeMux) {
	rootMux.Handle("/assets/", SetupAssetsRoutes())
	rootMux.Handle("/", SetupHomeRoutes())
	rootMux.Handle("/api/", SetupApiRoutes())
	rootMux.Handle("/auth/", SetupAuthRoutes())
}
