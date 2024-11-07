package middleware

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/utils"
)

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
		fmt.Println(cookie.Value, count)
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
