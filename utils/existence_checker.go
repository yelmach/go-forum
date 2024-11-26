package utils

import "forum/database"

func ExistsPost(postId int) bool {
	isValid := true
	if err := database.DataBase.QueryRow("SELECT EXISTS (SELECT * FROM posts WHERE id = ?)", postId).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}

func ExistsComment(commentId int) bool {
	isValid := true
	if err := database.DataBase.QueryRow(`SELECT EXISTS (SELECT * FROM comments WHERE id = ?)`, commentId).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}
