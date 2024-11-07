package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/database"
	"forum/models"
	"net/http"
)

// StoreSession is designed to save a new user session in a database if it doesn't already exist
func StoreSession(w http.ResponseWriter, session_id string, user models.User) error {
	// check for already stored session
	var count int
	err := database.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE user_id = ? ", user.Id).Scan(&count)
	if err != nil {
		return err
	}

	query := ``
	switch {
	case count > 0:
		query := `UPDATE sessions SET session_id = ?, expired_at = ? WHERE user_id = ?`
		if _, err := database.DataBase.Exec(query, session_id, 5, user.Id); err != nil {
			return err
		}
		return nil
	case count == 0:
		query = `INSERT INTO sessions (user_id, session_id, expired_at) VALUES (?, ?, ?)`
		if _, err := database.DataBase.Exec(query, user.Id, session_id, 5); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func GetSession(r *http.Request) (models.User, error) {
	id := r.Header["Authorization"]
	// fmt.Println(r.Header)
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
