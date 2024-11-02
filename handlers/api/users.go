package api

import (
	"encoding/json"
	"net/http"
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
	if postContent.Title == "" || postContent.Content == "" {
		http.Error(w, "err.Error()", http.StatusBadRequest)
		return
	}
	err = controllers.CreatePost(postContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "./web/templates/create_posts.html") // Serve your HTML form
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
