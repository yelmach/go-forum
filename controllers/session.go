package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"forum/database"
	"forum/models"
)

// StoreSession is designed to save a new user session in a database if it doesn't already exist
func StoreSession(w http.ResponseWriter, session_id string, user models.User) (int, error) {
	// check for already stored session
	var count int
	err := database.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE user_id = ? ", user.Id).Scan(&count)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, errors.New("user not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error scanning row: %w", err)
	}

	query := ``
	switch {
	case count > 0:
		query := `UPDATE sessions SET session_id = ? WHERE user_id = ?`
		if _, err := database.DataBase.Exec(query, session_id, user.Id); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	case count == 0:
		query = `INSERT INTO sessions (user_id, session_id) VALUES (?, ?)`
		if _, err := database.DataBase.Exec(query, user.Id, session_id); err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}

	return http.StatusOK, nil
}

func GetSession(r *http.Request) (models.User, error) {
	id := r.Header["Authorization"]
	if len(id) != 1 {
		return models.User{}, errors.New("no session id provided")
	}
	// get the id and the user from the db
	var user models.User
	stmt, err := database.DataBase.Prepare("SELECT user_id FROM sessions WHERE session_id = ?")
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

	stmt, err = database.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE id = ?")
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
