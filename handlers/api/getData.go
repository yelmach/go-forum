package api

import (
	"forum/database"
	"forum/models"
	"log"
)

func getReaction(Id int, ispost bool) ([]int, []int, error) {
	var queryLikes, queryDislikes string

	switch ispost {
	case true:
		queryLikes = `SELECT user_id FROM reactions WHERE post_id=? AND is_like=1`
		queryDislikes = `SELECT user_id FROM reactions WHERE post_id=? AND is_like=0`
	case false:
		queryLikes = `SELECT user_id FROM reactions WHERE comment_id=? AND is_like=1`
		queryDislikes = `SELECT user_id FROM reactions WHERE comment_id=? AND is_like=0`
	}
	userlikes, err := getUsersIds(queryLikes, Id)
	if err != nil {
		return []int{}, []int{}, err
	}

	userdislikes, err := getUsersIds(queryDislikes, Id)
	if err != nil {
		return []int{}, []int{}, err
	}

	return userlikes, userdislikes, nil
}

func getUsersIds(query string, Id int) ([]int, error) {
	usersIds := []int{}
	rows, err := database.DataBase.Query(query, Id)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var userid int

		if err := rows.Scan(&userid); err != nil {
			return []int{}, err
		}
		usersIds = append(usersIds, userid)
	}

	return usersIds, nil
}

func getUsername(userId int) (string, error) {
	var username string
	query := `SELECT username FROM users WHERE id=?`
	err := database.DataBase.QueryRow(query, userId).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func getPostComments(postId int) ([]models.CommentsApi, error) {
	comments := []models.CommentsApi{}

	query := `SELECT id, user_id, content, created_at FROM comments WHERE post_id=? ORDER BY created_at DESC`
	dbComments, err := database.DataBase.Query(query, postId)
	if err != nil {
		return []models.CommentsApi{}, err
	}
	defer dbComments.Close()

	for dbComments.Next() {
		var comment models.CommentsApi
		var userId int
		var commentid int

		err := dbComments.Scan(&commentid, &userId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return []models.CommentsApi{}, err
		}

		comment.Likes, comment.Dislikes, err = getReaction(commentid, false)
		if err != nil {
			return []models.CommentsApi{}, err
		}

		comment.By, err = getUsername(userId)
		if err != nil {
			return []models.CommentsApi{}, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func getPostCategories(postId int) ([]string, error) {
	categories := []string{}

	query := `SELECT category_id FROM post_categories WHERE post_id=?`
	queryRow, err := database.DataBase.Query(query, postId)
	if err != nil {
		return []string{}, err
	}
	defer queryRow.Close()

	for queryRow.Next() {
		var category_id int
		var content string
		if err := queryRow.Scan(&category_id); err != nil {
			log.Fatal(err)
		}

		query = `SELECT name FROM categories WHERE id=?`
		err = database.DataBase.QueryRow(query, category_id).Scan(&content)
		if err != nil {
			return []string{}, err
		}
		categories = append(categories, content)
	}

	return categories, nil
}
