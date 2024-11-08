package utils

import (
	"encoding/json"
	"net/http"
)

type Resp struct {
	Msg       string `json:"msg,omitempty"`
	Code      int    `json:"code,omitempty"`
	SessionId string `json:"session_id,omitempty"`
}

func ResponseJSON(w http.ResponseWriter, resp Resp) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)
}

func AddCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		MaxAge: 60 * 60 * 24,
	})
}

func DeleteCookie(w http.ResponseWriter, session_name string) {
	http.SetCookie(w, &http.Cookie{
		Name:   session_name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
