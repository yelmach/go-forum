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

// RegisterUser handles regestration request
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "invalid credentials", Code: http.StatusBadRequest})
		return
	}

	// check if the data provided exists
	if user.Username == "" || user.Email == "" || user.Password == "" {
		utils.ResponseJSON(w, utils.Resp{Msg: "username, email and password are required", Code: http.StatusBadRequest})
		return
	}

	// check username
	valid_username, err := utils.CheckUsernameFormat(user.Username)
	if err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	} else if !valid_username {
		utils.ResponseJSON(w, utils.Resp{Msg: "Invalid username format", Code: http.StatusBadRequest})
		return
	}

	// check email
	valid_email, err := utils.CheckEmailFormat(user.Email)
	if err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	} else if !valid_email {
		utils.ResponseJSON(w, utils.Resp{Msg: "Invalid email format", Code: http.StatusBadRequest})
		return
	}

	// check password
	if !utils.CheckPasswordFormat(user.Password) {
		utils.ResponseJSON(w, utils.Resp{Msg: "Invalid password format", Code: http.StatusBadRequest})
		return
	}

	// check user if exist
	if err := utils.CheckUserExist(user); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	// store user in database
	if err := controllers.RegisterUser(user); err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	loginToForum(w, r, user)
}

// LoginUser it handles login request
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if len(user.Username) > 60 || len(user.Password) > 20 {
		utils.ResponseJSON(w, utils.Resp{Msg: "Username or Password Incorrect", Code: http.StatusBadRequest})
		return
	}

	loginToForum(w, r, user)
}

// loginToForum logged the user to forum and create a session for that user
func loginToForum(w http.ResponseWriter, r *http.Request, user models.User) {
	// check user if exists
	user, statuscode, err := controllers.LoginUser(user)
	if statuscode == http.StatusUnauthorized {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	} else if err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// create session
	id, err := uuid.NewV7()
	if err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	sessionId := id.String()

	// store session in database
	statuscode, err = controllers.StoreSession(w, sessionId, user)
	if err != nil {
		handlers.ErrorHandler(w, r, statuscode)
		return
	}

	utils.AddCookie(w, "session_id", sessionId)
	utils.AddCookie(w, "user_id", strconv.Itoa(user.Id))
	utils.AddCookie(w, "username", user.Username)

	utils.ResponseJSON(w, utils.Resp{Msg: "Logged in successfuly", Code: http.StatusOK, SessionId: sessionId})
}

// LogOutUser it handles log out request
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlers.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	session_id, err := r.Cookie("session_id")
	if err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	query := `DELETE FROM sessions WHERE session_id=?`
	if _, err := database.DataBase.Exec(query, session_id.Value); err != nil {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	utils.DeleteCookie(w, "session_id")
	utils.DeleteCookie(w, "user_id")
	utils.DeleteCookie(w, "username")

	http.Redirect(w, r, "/", http.StatusFound)
}
