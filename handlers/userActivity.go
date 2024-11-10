package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/controllers"
	"forum/utils"

	"forum/models"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	Cookie_user_id, err := r.Cookie("user_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}
	user_id, err := strconv.Atoi(Cookie_user_id.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}

	postContent := models.Post{
		UserId:     user_id,
		Title:      r.FormValue("Title"),
		Content:    r.FormValue("Content"),
		CategoryId: r.Form["categories"],
		ImageUrl:   r.FormValue("Image_url"),
	}

	if postContent.Title == "" || postContent.Content == "" {
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
	w.WriteHeader(http.StatusCreated)
	http.ServeFile(w, r, "./web/templates/create_posts.html")
}

func CreateCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	name_categorie := r.URL.Query().Get("categori_name")

	if len(name_categorie) == 0 {
		utils.ResponseJSON(w, utils.Resp{Msg: "categorie name should be provided", Code: http.StatusBadRequest})
		return
	}

	statuscode, err := controllers.CreateCategorie(name_categorie)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: statuscode})
		return
	}
}

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := models.Comment{}
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	comment.UserId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "bad request", Code: http.StatusBadRequest})
		return
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
}

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	reaction := models.Reaction{}
	err := json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	cookie, _ := r.Cookie("user_id")
	reaction.UserId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if reaction.CommentId == 0 && reaction.PostId == 0 {
		utils.ResponseJSON(w, utils.Resp{Msg: "Either post_id or comment_id must be provided", Code: http.StatusBadRequest})
		return
	}

	err = controllers.CreateReaction(reaction)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
