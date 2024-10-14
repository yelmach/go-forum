package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	ID       string `json:"id`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func insialationDb() {
	var err error
	db, err = sql.Open("sqlite3", "forim.db")
	if err != nil {
		log.Fatal("error opening", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
	id TEXT PRIMERY KEY,
	email TEXT UNIQUE,
	username TEXT UNIQUE,
	password TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func Checkregester(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "information not correct", http.StatusBadRequest)
		return
	}

	// pass := []byte("People talk things not reasonable but you need not worry")

	// Hashing the password
	hashpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(hash))
	// user.ID = uuid.New().string()

	if err != nil {
		http.Error(w, "COULDN'T HASHONG THE CODE", http.StatusInternalServerError)
		return
	}
	uui, err := uuid.NewV4()
	user.ID = uui.String()
	_, err = db.Exec("INSERT INTO users(id, email, username, password)VALUES(?, ?, ?, ?)", user.ID, user.Email, user.Username, hashpass)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			http.Error(w, "Email already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) // 409
	json.NewEncoder(w).Encode("User registered seccsefuly")
}

func main() {
}
