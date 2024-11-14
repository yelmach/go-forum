package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/utils"
)

type Tempaltes struct {
	Root     *template.Template
	Register *template.Template
	Login    *template.Template
}

var (
	TemplateError error
	templates     Tempaltes
)

func init() {
	templates.Root, TemplateError = template.ParseFiles(
		"./web/templates/index.html",
		"./web/templates/components/guest_navbar.html",
		"./web/templates/components/guest_sidebar.html",
		"./web/templates/components/logged_navbar.html",
		"./web/templates/components/logged_sidebar.html",
	)
	templates.Register, TemplateError = template.ParseFiles("./web/templates/register.html")
	templates.Login, TemplateError = template.ParseFiles("./web/templates/login.html")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if TemplateError != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	IsLoggedIn := false

	cookie, err := r.Cookie("session_id")

	if err != nil {
		IsLoggedIn = false
		err = templates.Root.ExecuteTemplate(w, "index.html", IsLoggedIn)
	} else {
		count := 0

		if err := database.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE session_id=?", cookie.Value).Scan(&count); err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}

		if count == 0 {
			utils.DeleteCookie(w, "session_id")
			utils.DeleteCookie(w, "user_id")
			utils.DeleteCookie(w, "username")
			IsLoggedIn = false
			if err = templates.Root.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
				utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
			}

		}

		IsLoggedIn = true
		err = templates.Root.ExecuteTemplate(w, "index.html", IsLoggedIn)
	}

	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if TemplateError != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	templates.Register.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if TemplateError != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	templates.Login.Execute(w, nil)
}
