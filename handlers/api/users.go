package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/controllers"
	"forum/models"
	"forum/tools"
	"forum/utils"

	"github.com/gofrs/uuid"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = controllers.CreateUser(user)
	if err != nil {
		w.WriteHeader(400)
		data, err := json.Marshal(struct {
			Msg string
		}{
			Msg: err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(data)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err = controllers.LoginUser(user)
	if err != nil {
		w.WriteHeader(400)
		data, err := json.Marshal(struct {
			Msg string
		}{
			Msg: err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(data)
	}

	// create session
	id, err := uuid.NewV7()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	sessionId := id.String()

	// store session in database
	err = controllers.StoreSession(w, sessionId, user)
	if err != nil {
		http.Error(w, "You already have a session", http.StatusBadRequest)
	}

	tools.AddCookie(w, "session_id", sessionId)
	tools.AddCookie(w, "user_id", strconv.Itoa(user.Id))
	tools.AddCookie(w, "username", user.Username)

	w.WriteHeader(200)
	data, err := json.Marshal(struct {
		Msg       string
		SessionId string
	}{
		Msg:       "Logged in",
		SessionId: sessionId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(data)
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	// get session from database
	_, err := controllers.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct {
		Msg string
	}{
		Msg: "You are logged in",
	})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session_id := r.Cookies()[0].Value
	query := `DELETE FROM sessions WHERE session_id=?`
	if _, err := utils.DataBase.Exec(query, session_id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tools.DeleteCookie(w, "session_id")
	tools.DeleteCookie(w, "user_id")
	tools.DeleteCookie(w, "username")

	http.Redirect(w, r, "/", http.StatusFound)
}
