package routers

import (
	"net/http"

	"forum/handlers"
)

func SetupAssetsRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/assets/", handlers.AssetsHandler)

	return mux
}
