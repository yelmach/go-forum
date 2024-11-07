package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/database"
	"forum/models"
)

func LoadPostData(w http.ResponseWriter, r *http.Request) {
	var post models.PostsApi
	var userId int

	id, _ := strconv.Atoi(r.PathValue("id"))
	query := `SELECT id, user_id, title, content, image_url, created_at FROM posts WHERE id=?`
	err := database.DataBase.QueryRow(query, id).Scan(&post.Id, &userId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		newError := models.Error{}
		newError.Error.Status = http.StatusNotFound
		newError.Error.Code = "not_found"
		json.NewEncoder(w).Encode(newError)
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

	post.Likes, post.Dislikes, err = getReaction(post.Id, true)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func LoadData(w http.ResponseWriter, r *http.Request) {
	// Query to get all posts from the posts table
	query := `SELECT id, user_id, title, content, image_url, created_at FROM posts ORDER BY created_at DESC`

	// Execute the query
	dbPosts, err := database.DataBase.Query(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Prepare a slice to store the posts
	posts := []models.PostsApi{}

	// Iterate through the rows and scan each row into a Post struct
	for dbPosts.Next() {
		var post models.PostsApi
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

		post.Likes, post.Dislikes, err = getReaction(post.Id, true)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
	}

	// Check for errors from iterating over rows
	if err = dbPosts.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(posts)
}

func LoadAllCategories(w http.ResponseWriter, r *http.Request) {
	dbCategories, err := database.DataBase.Query(`SELECT id,name FROM categories`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbCategories.Close()

	categories := struct {
		Id   []int
		Name []string
	}{}
	for dbCategories.Next() {
		var category string
		var id int
		dbCategories.Scan(&id, &category)
		categories.Id = append(categories.Id, id)
		categories.Name = append(categories.Name, category)
	}
	
	if err = dbCategories.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categories)
}
