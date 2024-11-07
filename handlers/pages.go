package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/tools"
	"forum/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	IsLoggedIn := false

	tmpl, err := template.ParseFiles("./web/templates/index.html",
		"./web/templates/components/guest_navbar.html",
		"./web/templates/components/guest_sidebar.html",
		"./web/templates/components/logged_navbar.html",
		"./web/templates/components/logged_sidebar.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cookie, err := r.Cookie("session_id")

	if err != nil {
		IsLoggedIn = false
		err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn)
	} else {
		count := 0

		if err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE session_id=?", cookie.Value).Scan(&count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if count == 0 {
			tools.DeleteCookie(w, "session_id")
			tools.DeleteCookie(w, "user_id")
			tools.DeleteCookie(w, "username")
			IsLoggedIn = false
			if err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
				http.Error(w, "failled to execute temp", http.StatusInternalServerError)
			}

		}

		IsLoggedIn = true
		err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn)
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failled to execute template", http.StatusInternalServerError)
	}
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

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/create_posts.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}
