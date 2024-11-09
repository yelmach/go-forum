package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/database"
	"forum/utils"
)

func GetCookiesId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		is_valid := false
		cookie_session, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "not allowed", http.StatusBadRequest)
			return
		}
		session_id := cookie_session.Value
		cookie_user, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "not allowed", http.StatusBadRequest)
			return
		}
		user_id, err := strconv.Atoi(cookie_user.Value)
		fmt.Println("user_id :", user_id)
		if err != nil {
			http.Error(w, "not allowed", http.StatusBadRequest)
			return
		}
		if err := database.DataBase.QueryRow(`SELECT EXISTS(SELECT * FROM sessions JOIN users ON sessions.user_id = users.id WHERE session_id = ? AND user_id = ?)`, session_id, user_id).Scan(&is_valid); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if !is_valid {
			utils.DeleteCookie(w, "session_id")
			utils.DeleteCookie(w, "user_id")
			utils.DeleteCookie(w, "username")
			http.Redirect(w, r, "/", http.StatusFound)
		}
		next.ServeHTTP(w, r)
	}
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "not allowed", http.StatusBadRequest)
			return
		}
		count := 0
		if err := database.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE session_id=?", cookie.Value).Scan(&count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if count == 0 {
			utils.DeleteCookie(w, "session_id")
			utils.DeleteCookie(w, "user_id")
			utils.DeleteCookie(w, "username")

			http.Redirect(w, r, "/", http.StatusFound)
		}
		next.ServeHTTP(w, r)
	}
}

func RedirectMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
