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

// ResponseJSON responsible for responses
func ResponseJSON(w http.ResponseWriter, resp Resp) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)
}
