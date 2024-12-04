package routers

import (
	"net/http"

	"forum/handlers"
	"forum/handlers/api"
)

// Setup Api Routes
func SetupApiRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "api/posts" {
			handlers.ErrorHandler(w, r, http.StatusNotFound)
			return
		}
	})

	mux.HandleFunc("/api/posts", api.LoadData)
	mux.HandleFunc("/api/posts/{id}", api.LoadPostData)
	mux.HandleFunc("/api/categories", api.LoadAllCategories)
	return mux
}
