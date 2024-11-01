package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/models"
	"forum/utils"
)

func LoadData(w http.ResponseWriter, r *http.Request) {
	// Query to get all posts from the posts table
	query := `SELECT id, user_id, title, content, image_url, created_at FROM posts`

	// Execute the query
	dbPosts, err := utils.DataBase.Query(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Prepare a slice to store the posts
	posts := []models.Post{}

	// Iterate through the rows and scan each row into a Post struct
	for dbPosts.Next() {
		var post models.Post
		var userId int
		err := dbPosts.Scan(&post.Id, &userId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		post.By, err = getUsername(userId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		post.Comments, err = getPostComments(post.Id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		post.Categories, err = getPostCategories(post.Id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
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

func getUsername(userId int) (string, error) {
	var username string
	query := `SELECT username FROM users WHERE id=?`
	err := utils.DataBase.QueryRow(query, userId).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func getPostComments(postId int) ([]models.Comment, error) {
	comments := []models.Comment{}

	query := `SELECT user_id, content, created_at FROM comments WHERE post_id=?`
	dbComments, err := utils.DataBase.Query(query, postId)
	if err != nil {
		return []models.Comment{}, err
	}
	defer dbComments.Close()

	for dbComments.Next() {
		var comment models.Comment
		var userId int

		err := dbComments.Scan(&userId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return []models.Comment{}, err
		}

		comment.By, err = getUsername(userId)
		if err != nil {
			return []models.Comment{}, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func getPostCategories(postId int) ([]string, error) {
	categories := []string{}

	query := `SELECT category_id FROM post_categories WHERE post_id=?`
	queryRow, err := utils.DataBase.Query(query, postId)
	if err != nil {
		return []string{}, err
	}
	defer queryRow.Close()

	for queryRow.Next() {
		var category_id int
		var content string
		if err := queryRow.Scan(&category_id); err != nil {
			log.Fatal(err)
		}

		query = `SELECT name FROM categories WHERE id=?`
		err = utils.DataBase.QueryRow(query, category_id).Scan(&content)
		if err != nil {
			return []string{}, err
		}
		categories = append(categories, content)
	}

	return categories, nil
}
