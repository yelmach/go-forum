package controllers

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/database"
	"forum/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user models.User) (int, error) {
	// Create if the user exist
	var count int
	err := database.DataBase.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ? ", user.Email, user.Username).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("user already exist")
	}

	cryptedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	stmt, err := database.DataBase.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, string(cryptedPass))
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func LoginUser(user models.User) (models.User, error) {
	existUser := models.User{}
	stmt, err := database.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE username = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close() // Ensure statement is closed

	err = stmt.QueryRow(user.Username).Scan(&existUser.Id, &existUser.Username, &existUser.Email, &existUser.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, fmt.Errorf("error scanning row: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password))
	if err != nil {
		return models.User{}, err
	}

	return existUser, nil
}
