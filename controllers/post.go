package controllers

import (
	"fmt"

	"forum/database"
	"forum/models"
)

func CreateCategories() {
	var isValid bool
	if err := database.DataBase.QueryRow(`SELECT EXISTS(SELECT * FROM categories WHERE name = 'Technology')`).Scan(&isValid); err != nil {
		fmt.Printf("Error checking existence of category: %v\n", err)
	}
	if !isValid {
		if _, err := database.DataBase.Exec(`INSERT INTO categories(name) VALUES ('Technology'),('GoLang'),('Gaming'),('Sports'),('Programming'),('Zone01'),('Back-end'),('Front-end')`); err != nil {
			fmt.Printf("Error inserting categories: %v\n", err)
		}
	}
}

// CreatePost stores content of the post, and relation between posts and categories.
func CreatePost(post models.Post) error {
	res, err := database.DataBase.Exec(`INSERT INTO posts(user_id, title, content) VALUES(?, ?, ?)`, post.UserId, post.Title, post.Content)
	if err != nil {
		return err
	}

	postId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, category := range post.Categories {
		var categoryId int
		if err := database.DataBase.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryId); err != nil {
			return err
		}
		if _, err := database.DataBase.Exec(`INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)`, postId, categoryId); err != nil {
			return fmt.Errorf("error linking post to a category %s: %w", category, err)
		}
	}
	return nil
}

// CreateComment stores the comment in the database
func CreateComment(comment models.Comment) error {
	if _, err := database.DataBase.Exec(`INSERT INTO comments(post_id, user_id, content) VALUES(?, ?, ?)`, comment.PostId, comment.UserId, comment.Content); err != nil {
		return err
	}
	return nil
}

// CreateReaction stores the reaction in the database
func CreateReaction(r models.Reaction) error {
	liked, disliked := false, false
	isPost := r.PostId != 0

	if isPost {
		if err := database.DataBase.QueryRow(`SELECT EXISTS(SELECT is_like FROM reactions WHERE user_id=? AND post_id=? AND is_like=1)`, r.UserId, r.PostId).Scan(&liked); err != nil {
			return err
		}
		if err := database.DataBase.QueryRow(`SELECT EXISTS(SELECT is_like FROM reactions WHERE user_id=? AND post_id=? AND is_like=0)`, r.UserId, r.PostId).Scan(&disliked); err != nil {
			return err
		}

		if liked || disliked {
			if _, err := database.DataBase.Exec(`DELETE FROM reactions WHERE user_id=? AND post_id=?`, r.UserId, r.PostId); err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
		}

		if liked != r.IsLike || disliked != r.IsDislike {
			if _, err := database.DataBase.Exec(`INSERT INTO reactions (user_id, post_id, is_like) VALUES (?, ?, ?)`, r.UserId, r.PostId, r.IsLike); err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
		}
	} else {
		if err := database.DataBase.QueryRow(`SELECT EXISTS(SELECT is_like FROM reactions WHERE user_id=? AND comment_id=? AND is_like=1)`, r.UserId, r.CommentId).Scan(&liked); err != nil {
			return err
		}
		if err := database.DataBase.QueryRow(`SELECT EXISTS(SELECT is_like FROM reactions WHERE user_id=? AND comment_id=? AND is_like=0)`, r.UserId, r.CommentId).Scan(&disliked); err != nil {
			return err
		}

		if liked || disliked {
			if _, err := database.DataBase.Exec(`DELETE FROM reactions WHERE user_id=? AND comment_id=?`, r.UserId, r.CommentId); err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
		}

		if liked != r.IsLike || disliked != r.IsDislike {
			if _, err := database.DataBase.Exec(`INSERT INTO reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)`, r.UserId, r.CommentId, r.IsLike); err != nil {
				return fmt.Errorf("error executing statement: %w", err)
			}
		}
	}
	return nil
}
