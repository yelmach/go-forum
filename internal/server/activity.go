package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Struct to hold form data
type Credentials struct {
	Username string
	Password string
}

// Database connection
var db *sql.DB

// Initialize the SQLite3 database and create a user table if it doesn't exist
func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	// Create table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`
	if _, err = db.Exec(createTable); err != nil {
		log.Fatal(err)
	}
}

// Handler for serving the login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	tmp, err := template.ParseFiles("../web/templates/login.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err = tmp.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
	tmp.Execute(w, nil)
}

// Handler for processing login form submission
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if credentials are valid
		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Verify password use password hashing
		if password != storedPassword {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Login successful! Welcome, %s", username)
	}
}

func Sign_UpHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	if r.URL.Path != "/sign_up" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	tmp, err := template.ParseFiles("../web/templates/sign_up.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err = tmp.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
	tmp.Execute(w, nil)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	if r.URL.Path != "/sign_up" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	Validation(db, username, password)
	fmt.Fprintf(w, "Add User successful! Welcome, %s", username, password, email)
}

func Validation(db *sql.DB, username, password string) (bool, error) {
	var isuser string
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&isuser)
	if err != nil {
		return false, fmt.Errorf("This username is already in use. Try another name.")
	}
	if isuser == "" {
		insertUser(db, username, password)
		fmt.Println("successful")
		return true, nil
	} else {
		return false, fmt.Errorf("This username is already in use. Try another name.")
	}
}

func insertUser(db *sql.DB, username, password string) {
	reqiteINSERT := `INSERT INTO student(username, password) VALUES (?, ?)`
	statement, err := db.Prepare(reqiteINSERT)
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(username, password)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	// index
	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/index" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	tmp, err := template.ParseFiles("../web/templates/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err = tmp.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
	tmp.Execute(w, nil)
}
