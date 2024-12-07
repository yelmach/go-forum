package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DataBase *sql.DB

// InitDb opens an SQLITE database and creates tables
func InitDb() error {
	// Initialize database connection
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	query, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(query)); err != nil {
		return err
	}

	DataBase = db
	// // Add test comments for pagination testing
	// if err := addTestComments(); err != nil {
	// 	return err
	// }
	return err
}

// addTestComments adds 100 test comments to the first post
func addTestComments() error {
	// Get the first post ID
	var postId int
	err := DataBase.QueryRow("SELECT id FROM posts ORDER BY id DESC LIMIT 1").Scan(&postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No posts exist yet
		}
		return err
	}

	// Get the first user ID
	var userId int
	err = DataBase.QueryRow("SELECT id FROM users ORDER BY id DESC LIMIT 1").Scan(&userId)
	if err != nil {
		return err
	}

	// Start a transaction for better performance
	tx, err := DataBase.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Add 100 test comments
	stmt, err := tx.Prepare(`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 1; i <= 100; i++ {
		content := fmt.Sprintf("Test comment number %d for pagination testing", i)
		if _, err := stmt.Exec(postId, userId, content); err != nil {
			return err
		}
	}

	return tx.Commit()
}
