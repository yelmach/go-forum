package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/database"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser insert user information to user table
func RegisterUser(user models.User) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// insert data
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	if _, err := database.DataBase.Exec(query, user.Username, user.Email, string(hashedPass)); err != nil {
		return err
	}

	return nil
}

// LoginUser checks if the user info are exists, and correct in user table database
func LoginUser(user models.User) (models.User, int, error) {
	existUser := models.User{}
	// check if username already exist
	query := "SELECT id, username, email, password FROM users WHERE username = ? OR email = ?"
	err := database.DataBase.QueryRow(query, user.Username, user.Username).Scan(&existUser.Id, &existUser.Username, &existUser.Email, &existUser.Password)
	if err == sql.ErrNoRows {
		return models.User{}, http.StatusUnauthorized, fmt.Errorf("user not found")
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password)); err != nil {
		return models.User{}, http.StatusUnauthorized, fmt.Errorf("incorrect password")
	}
	return existUser, http.StatusOK, nil
}

// StoreSession is designed to save a new user session in a database if it doesn't already exist
func StoreSession(w http.ResponseWriter, session_id string, user models.User) (int, error) {
	// check session if already stored
	var isValid bool
	err := database.DataBase.QueryRow("SELECT EXISTS(SELECT * FROM sessions WHERE user_id = ?)", user.Id).Scan(&isValid)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	switch {
	case isValid:
		query := `UPDATE sessions SET session_id = ? WHERE user_id = ?`
		if _, err := database.DataBase.Exec(query, session_id, user.Id); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	case !isValid:
		query := `INSERT INTO sessions (user_id, session_id) VALUES (?, ?)`
		if _, err := database.DataBase.Exec(query, user.Id, session_id); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}

	return http.StatusOK, nil
}
