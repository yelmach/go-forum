package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Database connection
var db *sql.DB

// Initialize the SQLite3 database and create a user table if it doesn't exist
func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	// Create table if not exists
	// if _, err = db.Exec(CreateTableUsers, CreateTablePost, CreateTableEngagement); err != nil {
	if _, err = db.Exec(CreateTableUsers); err != nil {
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
		err := db.QueryRow(LoginQuery, username, username).Scan(&storedPassword)
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
	if r.URL.Path != "/sign_up" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		tmp, err := template.ParseFiles("../web/templates/sign_up.html")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		if err = tmp.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
		tmp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		AddUser(w, r)
	} else {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("r.Method  : ", r.Method)
}

func Index(w http.ResponseWriter, r *http.Request) {
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
