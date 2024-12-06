package handlers

import (
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	var DetailsError struct {
		Title  string
		Status int
		Method string
		Path   string
	}

	switch status {
	case http.StatusNotFound:
		DetailsError.Title = "Error 404 (Not Found)!!"
	case http.StatusMethodNotAllowed:
		DetailsError.Title = "Error 405 (Method Not Allowed)!!"
	case http.StatusInternalServerError:
		DetailsError.Title = "Error 500 (Internal Server Error)!!"
	case http.StatusBadRequest:
		DetailsError.Title = "Error 400 (Status Bad Request)"
	case http.StatusBadGateway:
		DetailsError.Title = "Error 502 (Status Bad Gateway)"
	}

	w.WriteHeader(status)
	DetailsError.Status = status
	DetailsError.Method = r.Method
	DetailsError.Path = r.URL.Path
	if err := templates.ExecuteTemplate(w, "error.html", DetailsError); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}
