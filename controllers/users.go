package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/models"
	"forum/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.User) (int, error) {
	// Create if the user exist
	var count int
	err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ? ", user.Email, user.Username).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("user already exist")
	}

	cryptedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	stmt, err := utils.DataBase.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, string(cryptedPass))
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func LoginUser(user models.User) (models.User, error) {
	existUser := models.User{}
	stmt, err := utils.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE username = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close() // Ensure statement is closed

	err = stmt.QueryRow(user.Username).Scan(&existUser.Id, &existUser.Username, &existUser.Email, &existUser.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, fmt.Errorf("error scanning row: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password))
	if err != nil {
		return models.User{}, err
	}

	return existUser, nil
}

// StoreSession is designed to save a new user session in a database if it doesn't already exist
func StoreSession(id string, user models.User) error {
	// check for already stored session
	var count int
	err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM sessions WHERE user_id = ? ", id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("session already exist")
	}

	// expiration time for this session

	stmt, err := utils.DataBase.Prepare("INSERT INTO sessions (user_id, session_id) VALUES (?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Id, id)
	if err != nil {
		return err
	}
	return nil
}

func GetSession(r *http.Request) (models.User, error) {
	id := r.Header["Authorization"]
	fmt.Println(r.Header)
	if len(id) != 1 {
		return models.User{}, errors.New("no session id provided")
	}
	// get the id and the user from the db
	var user models.User
	stmt, err := utils.DataBase.Prepare("SELECT user_id FROM sessions WHERE session_id = ?")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	var user_id int
	err = stmt.QueryRow(id[0]).Scan(&user_id)
	if err != nil {
		return models.User{}, err
	}
	fmt.Println(user_id)

	stmt, err = utils.DataBase.Prepare("SELECT id, username, email, password FROM users WHERE id = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close() // Ensure statement is closed

	err = stmt.QueryRow(user_id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}

func CreatePost(postContent models.PostContent) error {
	C_post, err := utils.DataBase.Prepare(`INSERT INTO posts(user_id, title, content, image_url, created_at) VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer C_post.Close()

	res, err := C_post.Exec(postContent.User_id, postContent.Title, postContent.Content, postContent.Image_url, postContent.Created_at)
	if err != nil {
		return err
	}

	// Get the ID of the newly created post
	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	C_postCategories, err := utils.DataBase.Prepare(`INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement for post_categories: %w", err)
	}
	defer C_postCategories.Close()

	for _, categoryID := range postContent.Category_id {
		Category_id, _ := strconv.Atoi(categoryID)
		if _, err := C_postCategories.Exec(postID, Category_id); err != nil {
			return fmt.Errorf("error linking post to category %s: %w", categoryID, err)
		}
	}
	return nil
}

func CreateComments(commentContent models.Comments) error {
	C_comment, err := utils.DataBase.Prepare(`INSERT INTO comments(post_id, user_id, content, created_at) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}

	defer C_comment.Close()

	if _, err = C_comment.Exec(commentContent.Post_id, commentContent.User_id, commentContent.Content, commentContent.Created_at); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateReaction(reactions models.Reactions) error {
	var query string
	var id int
	isPost := reactions.Post_id != 0
	if isPost {
		query = `INSERT INTO reactions (user_id, post_id, is_like, created_at) VALUES (?, ?, ?, ?)`
		id = reactions.Post_id
	} else {
		query = `INSERT INTO reactions (user_id, comment_id, is_like, created_at) VALUES (?, ?, ?, ?)`
		id = reactions.Comment_id
	}
	C_reaction, err := utils.DataBase.Prepare(query)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}
	if _, err = C_reaction.Exec(reactions.User_id, id, reactions.Is_like, reactions.Created_at); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}
