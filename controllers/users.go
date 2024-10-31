package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"forum/models"
	"forum/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.User) (int, error) {
	// Create if the user exist
	var count int
	err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ? ", user.Email, user.Username).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("user already exist")
	}

	cryptedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	stmt, err := utils.DataBase.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
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
	stmt, err := utils.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE username = ?")
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

// StoreSession is designed to save a new user session in a database if it doesn't already exist
func StoreSession(id string, user models.User) error {
	// check for already stored session
	var count int
	err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE user_id = ? ", id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("session already exist")
	}

	stmt, err := utils.DataBase.Prepare("INSERT INTO sessions (user_id, session_id) VALUES (?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Id, id)
	if err != nil {
		return err
	}
	return nil
}

func GetSession(r *http.Request) (models.User, error) {
	id := r.Header["Authorization"]

	if len(id) != 1 {
		return models.User{}, errors.New("no session id provided")
	}

	// get the id and the user from the db
	var user models.User
	stmt, err := utils.DataBase.Prepare("SELECT user_id FROM sessions WHERE session_id = ?")
	if err != nil {
		return models.User{}, err
	}

	defer stmt.Close()
	var user_id int
	err = stmt.QueryRow(id[0]).Scan(&user_id)
	if err != nil {
		return models.User{}, err
	}
	fmt.Println(user_id)

	stmt, err = utils.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE id = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close() // Ensure statement is closed

	err = stmt.QueryRow(user_id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}
