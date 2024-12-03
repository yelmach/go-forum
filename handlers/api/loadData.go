package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"forum/database"
	"forum/handlers"
	"forum/models"
	"forum/utils"
)

// LoadPostData gets data of one post from database and send it to js
func LoadPostData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	var post models.PostApi
	var userId int
	statuscode := 0

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "bad request", Code: http.StatusBadRequest})
		return
	}

	query := `SELECT id, user_id, title, content, image_url, created_at FROM posts WHERE id=?`
	err = database.DataBase.QueryRow(query, id).Scan(&post.Id, &userId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		newError := models.Error{}
		newError.Error.Status = http.StatusNotFound
		newError.Error.Code = "not_found"
		json.NewEncoder(w).Encode(newError)
		return
	}

	post.By, statuscode, err = getUsername(userId)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	post.Comments, statuscode, err = getPostComments(post.Id)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	post.Categories, statuscode, err = getPostCategories(post.Id)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	post.Likes, post.Dislikes, statuscode, err = getReaction(post.Id, true)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func LoadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// Get page parameter from query string, default is 1
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get filter type
	filterType := r.URL.Query().Get("filter")
	userID := r.URL.Query().Get("userId")
	category := r.URL.Query().Get("category")

	// Set posts per page
	const postsPerPage = 50
	offset := (page - 1) * postsPerPage

	var query string
	var args []interface{}

	switch filterType {
	case "created":
		query = `SELECT id, user_id, title, content, image_url, created_at 
                 FROM posts 
                 WHERE user_id = ? 
                 ORDER BY created_at DESC 
                 LIMIT ? OFFSET ?`
		args = []interface{}{userID, postsPerPage, offset}
	case "liked":
		query = `SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.created_at 
                 FROM posts p
                 JOIN reactions r ON p.id = r.post_id 
                 WHERE r.user_id = ? AND r.is_like = 1
                 ORDER BY p.created_at DESC 
                 LIMIT ? OFFSET ?`
		args = []interface{}{userID, postsPerPage, offset}
	case "category":
		query = `SELECT p.id, p.user_id, p.title, p.content, p.image_url, p.created_at 
                 FROM posts p
                 JOIN post_categories pc ON p.id = pc.post_id 
                 JOIN categories c ON pc.category_id = c.id
                 WHERE c.name = ?
                 ORDER BY p.created_at DESC 
                 LIMIT ? OFFSET ?`
		args = []interface{}{category, postsPerPage, offset}
	default:
		query = `SELECT id, user_id, title, content, image_url, created_at 
                 FROM posts 
                 ORDER BY created_at DESC 
                 LIMIT ? OFFSET ?`
		args = []interface{}{postsPerPage, offset}
	}

	dbPosts, err := database.DataBase.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, utils.Resp{Msg: "no rows found", Code: http.StatusNotFound})
			return
		}
		utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}
	defer dbPosts.Close()

	posts := []models.PostApi{}

	for dbPosts.Next() {
		var post models.PostApi
		var userId int
		var statuscode int

		if err := dbPosts.Scan(&post.Id, &userId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
			return
		}

		// Get all related data for the post
		post.By, statuscode, err = getUsername(userId)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
			return
		}

		post.Comments, statuscode, err = getPostComments(post.Id)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
			return
		}

		post.Categories, statuscode, err = getPostCategories(post.Id)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
			return
		}

		post.Likes, post.Dislikes, statuscode, err = getReaction(post.Id, true)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
			return
		}

		posts = append(posts, post)
	}

	// Get total count based on filter
	var countQuery string
	var countArgs []interface{}

	switch filterType {
	case "created":
		countQuery = "SELECT COUNT(*) FROM posts WHERE user_id = ?"
		countArgs = []interface{}{userID}
	case "liked":
		countQuery = `SELECT COUNT(*) FROM posts p
                      JOIN reactions r ON p.id = r.post_id 
                      WHERE r.user_id = ? AND r.is_like = 1`
		countArgs = []interface{}{userID}
	case "category":
		countQuery = `SELECT COUNT(*) FROM posts p
                      JOIN post_categories pc ON p.id = pc.post_id 
                      JOIN categories c ON pc.category_id = c.id
                      WHERE c.name = ?`
		countArgs = []interface{}{category}
	default:
		countQuery = "SELECT COUNT(*) FROM posts"
	}

	var totalPosts int
	if err := database.DataBase.QueryRow(countQuery, countArgs...).Scan(&totalPosts); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}

	response := struct {
		Posts      []models.PostApi `json:"posts"`
		TotalPosts int              `json:"totalPosts"`
		HasMore    bool             `json:"hasMore"`
		Page       int              `json:"page"`
	}{
		Posts:      posts,
		TotalPosts: totalPosts,
		HasMore:    offset+len(posts) < totalPosts,
		Page:       page,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// LoadAllCategories gets all categories from database and send it to js
func LoadAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	dbCategories, err := database.DataBase.Query(`SELECT name FROM categories`)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, utils.Resp{Msg: "no rows found", Code: http.StatusNotFound})
		} else {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		}
		return
	}
	defer dbCategories.Close()

	categories := []string{}

	for dbCategories.Next() {
		var category string
		if err := dbCategories.Scan(&category); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
			return
		}
		categories = append(categories, category)
	}

	if err = dbCategories.Err(); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}
}
