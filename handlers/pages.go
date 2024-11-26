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
	homeTemplateError     error
	registerTemplateError error
	loginTemplateError    error
	templates             Tempaltes
)

// parse all tamplates at once in the beggining of the program
func init() {
	templates.Root, homeTemplateError = template.ParseFiles(
		"./web/templates/index.html",
		"./web/templates/components/guest_navbar.html",
		"./web/templates/components/guest_sidebar.html",
		"./web/templates/components/logged_navbar.html",
		"./web/templates/components/logged_sidebar.html",
	)
	templates.Register, registerTemplateError = template.ParseFiles("./web/templates/register.html")
	templates.Login, loginTemplateError = template.ParseFiles("./web/templates/login.html")
}

// HomeHandler it handles requests to home page "/"
// execute the home page and show it to the user
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if homeTemplateError != nil {
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

	var IsLoggedIn bool

	// check if user already logged in from another browser
	cookie, err := r.Cookie("session_id")
	if err != nil {
		IsLoggedIn = false
		err = templates.Root.ExecuteTemplate(w, "index.html", IsLoggedIn)
	} else {
		var isValid bool
		if err := database.DataBase.QueryRow("SELECT EXISTS(SELECT * FROM sessions WHERE session_id=?)", cookie.Value).Scan(&isValid); err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}

		if !isValid {
			utils.DeleteCookie(w, "session_id")
			utils.DeleteCookie(w, "user_id")
			utils.DeleteCookie(w, "username")
			IsLoggedIn = false
			if err = templates.Root.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
		}

		IsLoggedIn = true
		err = templates.Root.ExecuteTemplate(w, "index.html", IsLoggedIn)
	}

	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// RegisterHandler it handles requests to register page "/register"
// parse the register page and show it to the user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if registerTemplateError != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if err := templates.Register.Execute(w, nil); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// RegisterHandler it handles requests to login page "/login"
// parse the login page and show it to the user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if loginTemplateError != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if err := templates.Login.Execute(w, nil); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
