package utils

import "forum/database"

func DelayPost() bool {
	isValid := false
	if err := database.DataBase.QueryRow(`SELECT EXISTS(
		SELECT * FROM posts JOIN users ON posts.user_id = users.id 
		WHERE created_at >= datetime('now', '-5 minutes'))`).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}

func DelayComment() bool {
	isValid := false
	if err := database.DataBase.QueryRow(`SELECT EXISTS(
		SELECT * FROM comments JOIN users ON comments.user_id = users.id 
		WHERE created_at >= datetime('now', '-20 seconds'))`).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}
