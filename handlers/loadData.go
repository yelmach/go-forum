package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/models"
	"forum/utils"
)

func LoadData(w http.ResponseWriter, r *http.Request) {
	// Query to get all posts from the posts table
	query := `SELECT id, user_id, title, content, image_url, created_at FROM posts`

	// Execute the query
	dbPosts, err := utils.DataBase.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPosts.Close()

	// Prepare a slice to store the posts
	posts := []models.Post{}

	// Iterate through the rows and scan each row into a Post struct
	for dbPosts.Next() {
		var post models.Post
		var posterId int
		err := dbPosts.Scan(&post.Id, &posterId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		query = `SELECT username FROM users WHERE id=?`
		err = utils.DataBase.QueryRow(query, posterId).Scan(&post.By)
		if err == sql.ErrNoRows {
			fmt.Fprintln(os.Stderr, "Error")
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "Error")
		}

		query = `SELECT user_id, content, created_at FROM comments WHERE post_id=?`
		dbComments, err := utils.DataBase.Query(query, post.Id)
		if err != nil {
			log.Fatal(err)
		}
		defer dbComments.Close()

		for dbComments.Next() {
			var comment models.Comment
			var commenterId int

			err := dbComments.Scan(&commenterId, &comment.Content, &comment.CreatedAt)
			if err != nil {
				log.Fatal(err)
			}

			query = `SELECT username FROM users WHERE id=?`
			err = utils.DataBase.QueryRow(query, commenterId).Scan(&comment.By)
			if err == sql.ErrNoRows {
				fmt.Fprintln(os.Stderr, "Error")
			} else if err != nil {
				fmt.Fprintln(os.Stderr, "Error")
			}
			post.Comments = append(post.Comments, comment)
		}

		posts = append(posts, post)
	}

	// Check for errors from iterating over rows
	if err = dbPosts.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(posts)
}
