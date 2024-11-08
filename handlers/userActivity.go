package handlers

import (
	"net/http"
	"strconv"

	"forum/controllers"
	"forum/utils"

	"forum/models"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Cookies()[1].Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
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
		utils.ResponseJSON(w, utils.Resp{Msg: "can't be empty", Code: http.StatusBadRequest})
		return
	} else if len(postContent.Title) >= 61 || len(postContent.Content) >= 2001 {
		utils.ResponseJSON(w, utils.Resp{Msg: "can't process, input to long", Code: http.StatusBadRequest})
		return
	}

	err = controllers.CreatePost(postContent)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	http.ServeFile(w, r, "./web/templates/create_posts.html")
}

func CreateCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	name_categorie := r.URL.Query().Get("categori_name")

	if len(name_categorie) == 0 {
		utils.ResponseJSON(w, utils.Resp{Msg: "categorie name should be provided", Code: http.StatusBadRequest})
		return
	}

	err, statuscode := controllers.CreateCategorie(name_categorie)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}
}

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed) // reprocess after
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	user_id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "bad request", Code: http.StatusBadRequest})
		return
	}

	post_id, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "bad request", Code: http.StatusBadRequest})
		return
	}

	comment := models.Comment{
		User_id: user_id,
		Post_id: post_id,
		Content: r.FormValue("content"),
	}

	if comment.Content == "" {
		utils.ResponseJSON(w, utils.Resp{Msg: "can't be empty", Code: http.StatusBadRequest})
		return
	} else if len(comment.Content) >= 501 {
		utils.ResponseJSON(w, utils.Resp{Msg: "cat't process, input to long", Code: http.StatusBadRequest})
		return
	}

	err = controllers.CreateComment(comment)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusCreated)
	HomeHandler(w, r)
}

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Cookies()[1].Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	postID := r.URL.Query().Get("post_id")
	commentID := r.URL.Query().Get("comment_id")
	isLike := r.URL.Query().Get("is_like")

	like, err := strconv.ParseBool(isLike)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	reactions := models.Reaction{
		User_id:    user_id,
		Post_id:    0,
		Comment_id: 0,
		Is_like:    like,
	}

	if postID != "" {
		reactions.Post_id, err = strconv.Atoi(postID)
	} else if commentID != "" {
		reactions.Comment_id, err = strconv.Atoi(commentID)
	} else {
		utils.ResponseJSON(w, utils.Resp{Msg: "Either post_id or comment_id must be provided", Code: http.StatusBadRequest})
		return
	}
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	err = controllers.CreateReaction(reactions)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusOK)
}
