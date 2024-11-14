package auth

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/controllers"
	"forum/database"
	"forum/handlers"
	"forum/models"
	"forum/utils"

	"github.com/gofrs/uuid"
)

// this func responsible for writing responses to me for debbuging

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	isValidPassword := utils.CheckPasswordFormat(user.Password)
	isValidEmail, err := utils.CheckEmailFormat(user.Email)
	if err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if isValidEmail && isValidPassword {
		if err := controllers.RegisterUser(user); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
			return
		}
	} else {
		utils.ResponseJSON(w, utils.Resp{Msg: "invalid format", Code: http.StatusBadRequest})
		return
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	user, statuscode, err := controllers.LoginUser(user)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	// create session
	id, err := uuid.NewV7()
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	sessionId := id.String()

	// store session in database
	statuscode, err = controllers.StoreSession(w, sessionId, user)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	utils.AddCookie(w, "session_id", sessionId)
	utils.AddCookie(w, "user_id", strconv.Itoa(user.Id))
	utils.AddCookie(w, "username", user.Username)

	utils.ResponseJSON(w, utils.Resp{Msg: "Logged in", Code: http.StatusOK, SessionId: sessionId})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session_id, err := r.Cookie("session_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	query := `DELETE FROM sessions WHERE session_id=?`
	if _, err := database.DataBase.Exec(query, session_id.Value); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	utils.DeleteCookie(w, "session_id")
	utils.DeleteCookie(w, "user_id")
	utils.DeleteCookie(w, "username")

	http.Redirect(w, r, "/", http.StatusFound)
}