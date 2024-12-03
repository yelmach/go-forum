package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/controllers"
	"forum/utils"

	"forum/models"
)

// NewPostHandler handles creation post request and store it to database
func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	post := models.Post{}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
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

	if utils.DelayPost(post.UserId) {
		utils.ResponseJSON(w, utils.Resp{Msg: "You can only post once every 5 minutes", Code: http.StatusBadRequest})
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

	// check if title and content are written and not too long
	if post.Title == "" || post.Content == "" {
		utils.ResponseJSON(w, utils.Resp{Msg: "input can't be empty", Code: http.StatusBadRequest})
		return
	} else if len(post.Title) >= 61 || len(post.Content) >= 2001 {
		utils.ResponseJSON(w, utils.Resp{Msg: "can't process, input too long", Code: http.StatusBadRequest})
		return
	}

	if err := controllers.CreatePost(post); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// NewCommentHandler handles the creation of new comment request
func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	comment := models.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	defer r.Body.Close()

	if !utils.ExistsPost(comment.PostId) {
		utils.ResponseJSON(w, utils.Resp{Msg: "bad request", Code: http.StatusBadRequest})
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}

	comment.UserId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}

	if utils.DelayComment(comment.PostId, comment.UserId) {
		utils.ResponseJSON(w, utils.Resp{Msg: "You can only post once every  20 seconds", Code: http.StatusBadRequest})
		return
	}

	// check if comment are written and not too long
	if comment.Content == "" {
		utils.ResponseJSON(w, utils.Resp{Msg: "comment can't be empty", Code: http.StatusBadRequest})
		return
	} else if len(comment.Content) >= 501 {
		utils.ResponseJSON(w, utils.Resp{Msg: "cat't process, comment too long", Code: http.StatusBadRequest})
		return
	}

	if err := controllers.CreateComment(comment); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// ReactionHandler handles the reaction request on a post or a comment
func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	reaction := models.Reaction{}
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if (reaction.PostId != 0 && reaction.CommentId != 0) || (reaction.PostId == 0 && reaction.CommentId == 0) {
		utils.ResponseJSON(w, utils.Resp{Msg: "Either post_id or comment_id must be provided", Code: http.StatusBadRequest})
		return
	}

	if !reaction.IsLike && !reaction.IsDislike {
		utils.ResponseJSON(w, utils.Resp{Msg: "new reaction must be provided", Code: http.StatusBadRequest})
		return
	}

	if reaction.PostId != 0 {
		if !utils.ExistsPost(reaction.PostId) {
			utils.ResponseJSON(w, utils.Resp{Msg: "Post not found", Code: http.StatusNotFound})
			return
		}
	} else {
		if !utils.ExistsComment(reaction.CommentId) {
			utils.ResponseJSON(w, utils.Resp{Msg: "Comment not found", Code: http.StatusNotFound})
			return
		}
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadGateway})
		return
	}

	reaction.UserId, err = strconv.Atoi(cookie.Value)
	if err != nil {
		utils.ResponseJSON(w, utils.Resp{Msg: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if err := controllers.CreateReaction(reaction); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
