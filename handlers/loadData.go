package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"forum/handlers/api"
	"forum/utils"
)

func LoadPostData(w http.ResponseWriter, r *http.Request) {
	var post api.Post
	var userId int

	id, _ := strconv.Atoi(r.PathValue("id"))
	query := `SELECT id, user_id, title, content, image_url, created_at FROM posts WHERE id=?`
	err := utils.DataBase.QueryRow(query, id).Scan(&post.Id, &userId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		newError := api.Error{}
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
	dbPosts, err := utils.DataBase.Query(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Prepare a slice to store the posts
	posts := []api.Post{}

	// Iterate through the rows and scan each row into a Post struct
	for dbPosts.Next() {
		var post api.Post
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
	dbCategories, err := utils.DataBase.Query(`SELECT name FROM categories`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbCategories.Close()

	categories := []string{}
	for dbCategories.Next() {
		var category string
		dbCategories.Scan(&category)
		categories = append(categories, category)
	}

	if err = dbCategories.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
		
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categories)
}

func getReaction(Id int, ispost bool) ([]int, []int, error) {
	var queryLikes, queryDislikes string

	switch ispost {
	case true:
		queryLikes = `SELECT user_id FROM reactions WHERE post_id=? AND is_like=1`
		queryDislikes = `SELECT user_id FROM reactions WHERE post_id=? AND is_like=0`
	case false:
		queryLikes = `SELECT user_id FROM reactions WHERE comment_id=? AND is_like=1`
		queryDislikes = `SELECT user_id FROM reactions WHERE comment_id=? AND is_like=0`
	}
	userlikes, err := getUsersIds(queryLikes, Id)
	if err != nil {
		return []int{}, []int{}, err
	}

	userdislikes, err := getUsersIds(queryDislikes, Id)
	if err != nil {
		return []int{}, []int{}, err
	}

	return userlikes, userdislikes, nil
}

func getUsersIds(query string, Id int) ([]int, error) {
	usersIds := []int{}
	rows, err := utils.DataBase.Query(query, Id)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var userid int

		if err := rows.Scan(&userid); err != nil {
			return []int{}, err
		}
		usersIds = append(usersIds, userid)
	}

	return usersIds, nil
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

func getPostComments(postId int) ([]api.Comment, error) {
	comments := []api.Comment{}

	query := `SELECT id, user_id, content, created_at FROM comments WHERE post_id=?`
	dbComments, err := utils.DataBase.Query(query, postId)
	if err != nil {
		return []api.Comment{}, err
	}
	defer dbComments.Close()

	for dbComments.Next() {
		var comment api.Comment
		var userId int
		var commentid int

		err := dbComments.Scan(&commentid, &userId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return []api.Comment{}, err
		}

		comment.Likes, comment.Dislikes, err = getReaction(commentid, false)
		if err != nil {
			return []api.Comment{}, err
		}

		comment.By, err = getUsername(userId)
		if err != nil {
			return []api.Comment{}, err
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
