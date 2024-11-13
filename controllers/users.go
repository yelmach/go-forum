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
	// check if data provided exists
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("email, username and password are required")
	}

	// check if user already registred
	var isExist bool
	if err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)",
		user.Email, user.Username).Scan(&isExist); err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("user already exist")
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// insert data
	if _, err := database.DataBase.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username, user.Email, string(hashedPass)); err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

// LoginUser checks if the user info are exists, and correct in user table database
func LoginUser(user models.User) (models.User, int, error) {
	existUser := models.User{}
	// check if username already exist
	err := database.DataBase.QueryRow("SELECT id, username, email, password FROM users WHERE username = ? OR email = ?", user.Username, user.Username).
		Scan(&existUser.Id, &existUser.Username, &existUser.Email, &existUser.Password)
	if err == sql.ErrNoRows {
		return models.User{}, http.StatusUnauthorized, fmt.Errorf("user not found")
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password))
	if err != nil {
		return models.User{}, http.StatusUnauthorized, err
	}
	return existUser, http.StatusOK, nil
}
