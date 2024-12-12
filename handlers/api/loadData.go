package api

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"forum/database"
	"forum/models"
	"forum/utils"
)

const POSTS_PER_PAGE = 10

// LoadPostData gets data of one post from database and send it to js
func LoadPostData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, utils.Resp{Msg: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}

	var (
		post       models.PostApi
		userId     int
		statuscode int
	)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "bad request", Code: http.StatusBadRequest})
		return
	}

	currentPage, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || currentPage < 1 {
		utils.ResponseJSON(w, utils.Resp{Msg: "Bad Request", Code: http.StatusBadRequest})
		return
	}

	query := `SELECT id, user_id, title, content, created_at FROM posts WHERE id=?`
	err = database.DataBase.QueryRow(query, id).Scan(&post.Id, &userId, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, utils.Resp{Msg: "No Rows Found", Code: http.StatusNotFound})
			return
		} else {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
			return
		}
	}

	post.By, statuscode, err = getUsername(userId)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}

	post.Comments, statuscode, post.TotalComments, post.HasMoreComments, err = getPostComments(post.Id, currentPage)
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

// LoadData gets all posts data from database and send it to js
func LoadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, utils.Resp{Msg: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}

	var (
		userId      int
		totalPosts  int
		countQuery  string
		filterQuery string
		countArgs   []interface{}
		filterArgs  []interface{}
	)

	// Get filter type
	filterBy := r.URL.Query().Get("filterBy")
	category := r.URL.Query().Get("category")

	if filterBy == "created" || filterBy == "liked" {
		cookie_token, err := r.Cookie("session_id")
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
			return
		}

		if err := database.DataBase.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", cookie_token.Value).Scan(&userId); err != nil {
			if err == sql.ErrNoRows {
				utils.ResponseJSON(w, utils.Resp{Msg: "unauthorized user", Code: http.StatusUnauthorized})
				return
			} else {
				utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
				return
			}
		}
	}

	switch filterBy {
	case "created":
		countArgs = append(countArgs, userId)
		filterArgs = append(filterArgs, userId)
		countQuery = "SELECT COUNT(*) FROM posts WHERE user_id = ?"
		filterQuery = `SELECT id, user_id, title, content, created_at FROM posts 
                 WHERE user_id = ? 
                 ORDER BY created_at DESC 
                 LIMIT ? OFFSET ?`
	case "liked":
		countArgs = append(countArgs, userId)
		filterArgs = append(filterArgs, userId)
		countQuery = `SELECT COUNT(*) FROM posts p
                 JOIN reactions r ON p.id = r.post_id 
                 WHERE r.user_id = ? AND r.is_like = 1`
		filterQuery = `SELECT p.id, p.user_id, p.title, p.content, p.created_at FROM posts p
                 JOIN reactions r ON p.id = r.post_id 
                 WHERE r.user_id = ? AND r.is_like = 1
                 ORDER BY p.created_at DESC 
                 LIMIT ? OFFSET ?`
	case "category":
		countArgs = append(countArgs, category)
		filterArgs = append(filterArgs, category)
		countQuery = `SELECT COUNT(*) FROM posts p
                 JOIN post_categories pc ON p.id = pc.post_id 
                 JOIN categories c ON pc.category_id = c.id
                 WHERE c.name = ?`
		filterQuery = `SELECT p.id, p.user_id, p.title, p.content, p.created_at FROM posts p
                 JOIN post_categories pc ON p.id = pc.post_id 
                 JOIN categories c ON pc.category_id = c.id
                 WHERE c.name = ?
                 ORDER BY p.created_at DESC 
                 LIMIT ? OFFSET ?`
	default:
		countQuery = "SELECT COUNT(*) FROM posts"
		filterQuery = `SELECT id, user_id, title, content, created_at FROM posts 
                 ORDER BY created_at DESC 
                 LIMIT ? OFFSET ?`
	}

	// Get total count based on filter
	if err := database.DataBase.QueryRow(countQuery, countArgs...).Scan(&totalPosts); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}

	currentPage, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || currentPage < 1 || float64(currentPage) > math.Ceil(float64(totalPosts)/float64(POSTS_PER_PAGE)) {
		utils.ResponseJSON(w, utils.Resp{Msg: "No Post Found", Code: http.StatusNotFound})
		return
	}

	// Set posts offset
	offset := (currentPage - 1) * POSTS_PER_PAGE
	filterArgs = append(filterArgs, POSTS_PER_PAGE, offset)

	// Get all filtered posts at the current page
	dbPosts, err := database.DataBase.Query(filterQuery, filterArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, utils.Resp{Msg: "No Post Found", Code: http.StatusNotFound})
			return
		} else {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
			return
		}
	}
	defer dbPosts.Close()

	posts := []models.PostApi{}

	for dbPosts.Next() {
		var (
			userId     int
			statuscode int
			post       models.PostApi
		)

		if err := dbPosts.Scan(&post.Id, &userId, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
			return
		}

		// Get all related data for the post
		post.By, statuscode, err = getUsername(userId)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
			return
		}

		if err := database.DataBase.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id=?`, post.Id).Scan(&post.TotalComments); err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
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

	response := struct {
		Posts       []models.PostApi `json:"posts"`
		TotalPosts  int              `json:"totalPosts"`
		HasMore     bool             `json:"hasMore"`
		CurrentPage int              `json:"currentPage"`
	}{
		Posts:       posts,
		TotalPosts:  totalPosts,
		HasMore:     offset+len(posts) < totalPosts,
		CurrentPage: currentPage,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// LoadAllCategories gets all categories from database and send it to js
func LoadAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, utils.Resp{Msg: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}

	dbCategories, err := database.DataBase.Query(`SELECT name FROM categories`)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, utils.Resp{Msg: "No Rows Found", Code: http.StatusNotFound})
			return
		} else {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
			return
		}
	}
	defer dbCategories.Close()

	categories := []string{}

	for dbCategories.Next() {
		var category string
		if err := dbCategories.Scan(&category); err != nil {
			if err == sql.ErrNoRows {
				utils.ResponseJSON(w, utils.Resp{Msg: "No Rows Found", Code: http.StatusNotFound})
				return
			} else {
				utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
				return
			}
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
