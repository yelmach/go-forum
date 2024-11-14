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

	dbPosts, err := database.DataBase.Query(`SELECT id, user_id, title, content, image_url, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ResponseJSON(w, utils.Resp{Msg: "no rows found", Code: http.StatusNotFound})
		} else {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		}
		return
	}
	defer dbPosts.Close()

	posts := []models.PostApi{}

	for dbPosts.Next() {
		var post models.PostApi
		var userId int
		var statuscode int

		err := dbPosts.Scan(&post.Id, &userId, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt)
		if err != nil {
			utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
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

		posts = append(posts, post)
	}

	if err = dbPosts.Err(); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(posts)
}

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
