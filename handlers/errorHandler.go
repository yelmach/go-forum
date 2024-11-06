package handlers

import (
	"html/template"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	tem, err := template.ParseFiles("./web/templates/error.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	var DetelsError struct {
		Title  string
		Status int
		Method string
		Path   string
	}

	switch status {
	case http.StatusNotFound:
		DetelsError.Title = "Error 404 (Not Found)!!"
	case http.StatusMethodNotAllowed:
		DetelsError.Title = "Error 405 (Method Not Allowed)!!"
	case http.StatusInternalServerError:
		DetelsError.Title = "Error 500 (Internal Server Error)!!"
	case http.StatusCreated:
		DetelsError.Title = "Error 201 (Status Created)!!"
	case http.StatusBadRequest:
		DetelsError.Title = "Error 400 (Status Bad Request)"
	case http.StatusBadGateway:
		DetelsError.Title = "Error 502 (Status Bad Gateway)"
	}
	
	DetelsError.Status = status
	DetelsError.Method = r.Method
	DetelsError.Path = r.URL.Path
	if err := tem.Execute(w, DetelsError); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return

	}
}
