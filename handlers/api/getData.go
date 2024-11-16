package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/database"
	"forum/models"
)

// getReaction gets user id of users that do like or dislike action
// from reactions table by (post/comment id)
func getReaction(Id int, isPost bool) ([]int, []int, int, error) {
	var queryLikes, queryDislikes string

	switch isPost {
	case true:
		queryLikes = `SELECT user_id FROM reactions WHERE post_id=? AND is_like=1`
		queryDislikes = `SELECT user_id FROM reactions WHERE post_id=? AND is_like=0`
	case false:
		queryLikes = `SELECT user_id FROM reactions WHERE comment_id=? AND is_like=1`
		queryDislikes = `SELECT user_id FROM reactions WHERE comment_id=? AND is_like=0`
	}

	userlikes, statuscode, err := getUsersIds(queryLikes, Id)
	if err != nil {
		return []int{}, []int{}, statuscode, err
	}

	userdislikes, statuscode, err := getUsersIds(queryDislikes, Id)
	if err != nil {
		return []int{}, []int{}, statuscode, err
	}

	return userlikes, userdislikes, http.StatusOK, nil
}

// getUsersIds gets user id of users that do like or dislike action on a post or comment
func getUsersIds(query string, Id int) ([]int, int, error) {
	usersIds := []int{}
	rows, err := database.DataBase.Query(query, Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []int{}, http.StatusNotFound, err
		} else {
			return []int{}, http.StatusInternalServerError, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var userid int

		if err := rows.Scan(&userid); err != nil {
			return []int{}, http.StatusInternalServerError, err
		}
		usersIds = append(usersIds, userid)
	}

	return usersIds, http.StatusOK, nil
}

// getUsername gets username from users table by user id
func getUsername(userId int) (string, int, error) {
	var username string

	query := `SELECT username FROM users WHERE id=?`
	if err := database.DataBase.QueryRow(query, userId).Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			return "", http.StatusNotFound, fmt.Errorf("no rows found: %v", err)
		} else {
			return "", http.StatusInternalServerError, fmt.Errorf("internal Server Error")
		}
	}
	return username, http.StatusOK, nil
}

// getPostComments gets all comments from comments table by post id
func getPostComments(postId int) ([]models.CommentApi, int, error) {
	comments := []models.CommentApi{}

	query := `SELECT id, user_id, content, created_at FROM comments WHERE post_id=? ORDER BY created_at DESC`
	dbComments, err := database.DataBase.Query(query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.CommentApi{}, http.StatusNotFound, err
		} else {
			return []models.CommentApi{}, http.StatusInternalServerError, err
		}
	}
	defer dbComments.Close()

	for dbComments.Next() {
		var comment models.CommentApi
		var userId int
		var statuscode int

		if err := dbComments.Scan(&comment.Id, &userId, &comment.Content, &comment.CreatedAt); err != nil {
			return []models.CommentApi{}, http.StatusInternalServerError, err
		}

		comment.Likes, comment.Dislikes, statuscode, err = getReaction(comment.Id, false)
		if err != nil {
			return []models.CommentApi{}, statuscode, err
		}

		comment.By, statuscode, err = getUsername(userId)
		if err != nil {
			return []models.CommentApi{}, statuscode, err
		}
		comments = append(comments, comment)
	}

	return comments, http.StatusOK, nil
}

// getPostCategories gets all categories that assosiated to a post by post id
func getPostCategories(postId int) ([]string, int, error) {
	categories := []string{}

	query := `SELECT category_id FROM post_categories WHERE post_id=?`
	queryRow, err := database.DataBase.Query(query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, http.StatusNotFound, err
		} else {
			return []string{}, http.StatusInternalServerError, err
		}
	}
	defer queryRow.Close()

	for queryRow.Next() {
		var category_id int
		var content string

		if err := queryRow.Scan(&category_id); err != nil {
			return []string{}, http.StatusInternalServerError, err
		}

		query = `SELECT name FROM categories WHERE id=?`
		err = database.DataBase.QueryRow(query, category_id).Scan(&content)
		if err != nil {
			if err == sql.ErrNoRows {
				return []string{}, http.StatusNotFound, err
			} else {
				return []string{}, http.StatusInternalServerError, err
			}
		}

		categories = append(categories, content)
	}

	return categories, http.StatusOK, err
}
