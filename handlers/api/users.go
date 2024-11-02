package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/controllers"
	"forum/models"

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
	err = controllers.StoreSession(sessionId, user)
	if err != nil {
		http.Error(w, "You already have a session", http.StatusBadRequest)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   1000,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(user.Id),
	})

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
