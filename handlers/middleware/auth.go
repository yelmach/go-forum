package middleware

import (
	"net/http"
	"strconv"

	"forum/database"
	"forum/utils"
)

// Middleware allows only users that authenticated to use next handler(add reaction, comment or post)
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		is_valid := false
		cookie_session, err := r.Cookie("session_id")
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
			return
		}
		session_id := cookie_session.Value
		cookie_user, err := r.Cookie("user_id")
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
			return
		}

		user_id, err := strconv.Atoi(cookie_user.Value)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
			return
		}

		cookie_username, err := r.Cookie("username")
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
			return
		}

		user_username := cookie_username.Value
		query := `SELECT EXISTS(SELECT * FROM sessions JOIN users ON sessions.user_id = users.id WHERE session_id = ? AND user_id = ? AND users.username = ? )`
		if err := database.DataBase.QueryRow(query, session_id, user_id, user_username).Scan(&is_valid); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "internal server error", Code: http.StatusInternalServerError})
			return
		}
		if !is_valid {
			utils.DeleteCookie(w, "session_id")
			utils.DeleteCookie(w, "user_id")
			utils.DeleteCookie(w, "username")
			utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
			return
		}
		next.ServeHTTP(w, r)
	}
}

// RedirectMiddleware redirect the logged user to home page, if he
// tries to reach login and register page
func RedirectMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
