package utils

import "forum/database"

func IspostId(postId int) bool {
	isValid := true
	if err := database.DataBase.QueryRow("SELECT EXISTS (SELECT * FROM posts WHERE id = ?)", postId).Scan(&isValid); err != nil || !isValid {
		return false
	}
	return true
}
