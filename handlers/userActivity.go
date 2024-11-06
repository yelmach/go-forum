package handlers

import (
	"net/http"
	"strconv"

	"forum/controllers"

	"forum/models"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Cookies()[1].Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	postContent := models.Post{
		User_id:     user_id,
		Title:       r.FormValue("Title"),
		Content:     r.FormValue("Content"),
		Category_id: r.Form["categories"],
		Image_url:   r.FormValue("Image_url"),
	}
	if postContent.Title == "" || postContent.Content == "" || len(postContent.Category_id) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = controllers.CreatePost(postContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "./web/templates/create_posts.html")
}

func CreateCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	name_categorie := r.URL.Query().Get("categori_name")
	if len(name_categorie) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := controllers.CreateCategorie(name_categorie); err != nil {
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return

	}
}

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, _ := r.Cookie("user_id")
	user_id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	post_id, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	comment := models.Comment{
		User_id: user_id,
		Post_id: post_id,
		Content: r.FormValue("content"),
	}

	if comment.Content == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = controllers.CreateComment(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	HomeHandler(w, r)
}

func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	// id, _ := strconv.Atoi(r.PathValue("id"))

	user_id, err := strconv.Atoi(r.Cookies()[1].Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	postID := r.URL.Query().Get("post_id")
	commentID := r.URL.Query().Get("comment_id")
	isLike := r.URL.Query().Get("is_like")

	like, err := strconv.ParseBool(isLike)
	if err != nil {
		http.Error(w, "Invalid like value", http.StatusBadRequest)
		return
	}
	reactions := models.Reaction{
		User_id:    user_id,
		Post_id:    0,
		Comment_id: 0,
		Is_like:    like,
	}

	if postID != "" {
		reactions.Post_id, _ = strconv.Atoi(postID)
	} else if commentID != "" {
		reactions.Comment_id, _ = strconv.Atoi(commentID)
	} else {
		http.Error(w, "Either post_id or comment_id must be provided", http.StatusBadRequest)
		return
	}

	err = controllers.CreateReaction(reactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
