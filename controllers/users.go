package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"forum/database"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user models.User) error {
	cryptedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	stmt, err := database.DataBase.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(strings.ToLower(user.Username), strings.ToLower(user.Email), string(cryptedPass))
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(user models.User) (models.User, int, error) {
	existUser := models.User{}
	stmt, err := database.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE username = ? OR email = ?")
	if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close() // Ensure statement is closed

	err = stmt.QueryRow(strings.ToLower(user.Username), strings.ToLower(user.Username)).Scan(&existUser.Id, &existUser.Username, &existUser.Email, &existUser.Password)
	if err == sql.ErrNoRows {
		return models.User{}, http.StatusNotFound, errors.New("user not found")
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password))
	if err != nil {
		return models.User{}, http.StatusUnauthorized, err
	}
	return existUser, http.StatusOK, nil
}
