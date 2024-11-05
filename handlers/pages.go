package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	IsLoggedIn := false

	tmpl, err := template.ParseFiles("./web/templates/index.html",
		"./web/templates/components/guest_navbar.html",
		"./web/templates/components/guest_sidebar.html",
		"./web/templates/components/logged_navbar.html",
		"./web/templates/components/logged_sidebar.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	_, err = r.Cookie("session_id")
	if err != nil {
		IsLoggedIn = false
		err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn)
	} else {
		IsLoggedIn = true
		err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn)
	}
	if err != nil {
		// fmt.Println(err)
		// http.Error(w, "failled to execute template", http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/login.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/create_posts.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
