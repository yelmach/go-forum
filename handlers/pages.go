package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/controllers"
	"forum/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Categories []models.Categories
		Posts      []models.Posts
	}
	tmpl, err := template.ParseFiles("./web/templates/index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	categories, err := controllers.DisplayCategories("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	posts, err := controllers.DisplayPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(posts)
	data := Data{
		Categories: categories,
		Posts:      posts,
	}
	tmpl.Execute(w, data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/login.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}
