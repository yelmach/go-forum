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
	post := models.Post{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}
	post.UserId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}

	// handle repeated category
	if !utils.HasUniqueCategories(post.Categories) {
		utils.ResponseJSON(w, utils.Resp{Msg: "repeted category", Code: http.StatusBadRequest})
		return
	}

	// and not exists categories
	if err = utils.VerifyCategoriesMatch(post.Categories); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: "category not found", Code: http.StatusBadRequest})
		return
	}
	//

	if post.Title == "" || post.Content == "" {
		utils.ResponseJSON(w, utils.Resp{Msg: "can't be empty", Code: http.StatusBadRequest})
		return
	} else if len(post.Title) >= 61 || len(post.Content) >= 2001 {
		utils.ResponseJSON(w, utils.Resp{Msg: "can't process, input to long", Code: http.StatusBadRequest})
		return
	}

	err = controllers.CreatePost(post)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(http.StatusCreated)
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
		utils.ResponseJSON(w, utils.Resp{Msg: "cat't process, comment to long", Code: http.StatusBadRequest})
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

	if !reaction.IsLike && !reaction.IsDislike {
		utils.ResponseJSON(w, utils.Resp{Msg: "new reaction must be provided", Code: http.StatusBadRequest})
		return
	}

	err = controllers.CreateReaction(reaction)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
