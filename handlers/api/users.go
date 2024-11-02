package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"forum/controllers"
	"forum/models"

	"github.com/gofrs/uuid"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = controllers.CreateUser(user)
	if err != nil {
		w.WriteHeader(400)
		data, err := json.Marshal(struct {
			Msg string
		}{
			Msg: err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(data)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err = controllers.LoginUser(user)
	if err != nil {
		w.WriteHeader(400)
		data, err := json.Marshal(struct {
			Msg string
		}{
			Msg: err.Error(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(data)
	}

	// create session
	id, err := uuid.NewV7()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	sessionId := id.String()

	// store session in database
	err = controllers.StoreSession(sessionId, user)
	if err != nil {
		http.Error(w, "You already have a session", http.StatusBadRequest)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   1000,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(user.Id),
	})

	w.WriteHeader(200)
	data, err := json.Marshal(struct {
		Msg       string
		SessionId string
	}{
		Msg:       "Logged in",
		SessionId: sessionId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(data)
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	// get session from database
	_, err := controllers.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct {
		Msg string
	}{
		Msg: "You are logged in",
	})
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := controllers.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postContent := models.PostContent{
		User_id:     user.Id,
		Title:       r.FormValue("Title"),
		Content:     r.FormValue("Content"),
		Category_id: r.Form["categories"],
		Image_url:   r.FormValue("Image_url"),
		Created_at:  time.Now().Format("2006-01-02 15:04:05"),
	}
	if postContent.Title == "" || postContent.Content == "" || len(postContent.Category_id) == 0 {
		http.Error(w, "err.Error()", http.StatusBadRequest)
		return
	}
	err = controllers.CreatePost(postContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "./web/templates/create_posts.html")
}

func CreateCommentsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := controllers.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commentContent := models.Comments{
		User_id:    user.Id,
		Post_id:    1,
		Content:    r.FormValue("Content"),
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
	}
	if commentContent.Content == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	err = controllers.CreateComments(commentContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddLikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := controllers.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	reactions := models.Reactions{
		User_id:    user.Id,
		Post_id:    0,
		Comment_id: 0,
		Is_like:    like,
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
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
