package handlers

import (
	"net/http"

	"forum/database"
	"forum/utils"
)

// HomeHandler it handles requests to home page "/"
// execute the home page and show it to the user
func HomeHandler(w http.ResponseWriter, r *http.Request) {
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
	cookie, errNoCookie := r.Cookie("session_id")
	if errNoCookie != nil {
		IsLoggedIn = false
		if err := templates.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
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
			if err := templates.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
				ErrorHandler(w, r, http.StatusInternalServerError)
				return
			}
		}

		IsLoggedIn = true
		if err := templates.ExecuteTemplate(w, "index.html", IsLoggedIn); err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}

// RegisterHandler it handles requests to register page "/register"
// parse the register page and show it to the user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if err := templates.ExecuteTemplate(w, "register.html", nil); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// RegisterHandler it handles requests to login page "/login"
// parse the login page and show it to the user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
