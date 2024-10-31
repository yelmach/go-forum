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

func DisplayPosts(posts models.Posts) (string, error) {
	var resalt string
	err := utils.DataBase.QueryRow("SELECT id,title,content FROM posts").Scan(&resalt)
	if err != nil {
		return "", err
	}
	return resalt, nil
}
