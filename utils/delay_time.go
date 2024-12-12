package utils

import "forum/database"

func DelayPost(UserId int) bool {
	isValid := false
	if err := database.DataBase.QueryRow(`SELECT EXISTS(
		SELECT * FROM posts
		WHERE created_at >= datetime('now', '-1 minutes') 
		AND user_id =? )`, UserId).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}

func DelayComment(PostId, UserId int) bool {
	isValid := false
	if err := database.DataBase.QueryRow(`SELECT EXISTS(
		SELECT * FROM comments
		WHERE created_at >= datetime('now', '-5 seconds') 
		AND post_id = ? 
		AND user_id = ?)`, PostId, UserId).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}
