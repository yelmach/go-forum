package controllers

import (
	"forum/models"
	"forum/utils"
)

func CreatePosts(posts models.Posts) (bool, error) {
	stmt, err := utils.DataBase.Prepare("INSERT INTO posts (title,content,user_id) VALUES(?,?,?)")
	if err != nil {
		return false, err
	}
	stmt.Close()
	_, err = stmt.Exec(posts.Content, posts.Title, posts.UserId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DisplayPosts1(posts models.Posts) (string, error) {
	var resalt string
	err := utils.DataBase.QueryRow("SELECT id,title,content FROM posts").Scan(&resalt)
	if err != nil {
		return "", err
	}
	return resalt, nil
}

func DisplayPosts() ([]models.Posts, error) {
	var posts []models.Posts
	query := "SELECT * FROM posts"

	rows, err := utils.DataBase.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Posts
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
