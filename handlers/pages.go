package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(`D:\github\forum-01\web\templates\index.html`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(`D:\github\forum-01\web\templates\register.html`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(`D:\github\forum-01\web\templates\login.html`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}

