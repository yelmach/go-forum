package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	Validation(db, username, password, email)
	fmt.Fprintf(w, "Add User successful! Welcome, %s", username, password, email)
}

func Validation(db *sql.DB, username, password, email string) {
	var isuser string
	err := db.QueryRow(ValidUserQuery, username).Scan(&isuser)
	if err != nil {
		fmt.Errorf("This username is already in use. Try another name.")
	}
	if isuser == "" {
		insertUser(db, username, password, email)
		fmt.Println("successful")
	} else {
		fmt.Errorf("This username is already in use. Try another name.")
	}
}

func insertUser(db *sql.DB, username, password, email string) {
	statement, err := db.Prepare(AddUserQuery)
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(username, password, email)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
