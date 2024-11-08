package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/database"
	"forum/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	IsLoggedIn := false

	tmpl, err := template.ParseFiles("./web/templates/index.html",
		"./web/templates/components/guest_navbar.html",
		"./web/templates/components/guest_sidebar.html",
		"./web/templates/components/logged_navbar.html",
		"./web/templates/components/logged_sidebar.html")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
	}

	cookie, err := r.Cookie("session_id")

	if err != nil {
		IsLoggedIn = false
		err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn)
	} else {
		count := 0

		if err := database.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE session_id=?", cookie.Value).Scan(&count); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		}

		if count == 0 {
			utils.DeleteCookie(w, "session_id")
			utils.DeleteCookie(w, "user_id")
			utils.DeleteCookie(w, "username")
			IsLoggedIn = false
			if err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
				utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
			}

		}

		IsLoggedIn = true
		err = tmpl.ExecuteTemplate(w, "index.html", IsLoggedIn)
	}

	if err != nil {
		fmt.Println(err)
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
	}

	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/login.html")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
	}

	tmpl.Execute(w, nil)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createpost" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/create_posts.html")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
	}

	tmpl.Execute(w, nil)
}
