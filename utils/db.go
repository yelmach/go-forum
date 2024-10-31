package utils

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DataBase *sql.DB

func InitDb() error {
	// Initialize database connection
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}
	query, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	if _, err = db.Exec(string(query)); err != nil {
		return err
	}

	DataBase = db

	return err
}
