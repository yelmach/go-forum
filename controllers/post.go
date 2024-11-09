package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/database"
	"forum/models"
)

func CreatePost(postContent models.Post) error {
	C_post, err := database.DataBase.Prepare(`INSERT INTO posts(user_id, title, content, image_url) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer C_post.Close()

	res, err := C_post.Exec(postContent.UserId, postContent.Title, postContent.Content, postContent.ImageUrl)
	if err != nil {
		return err
	}

	// Get the ID of the newly created post
	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	C_postCategories, err := database.DataBase.Prepare(`INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement for post_categories: %w", err)
	}
	defer C_postCategories.Close()

	for _, categoryID := range postContent.CategoryId {
		Category_id, _ := strconv.Atoi(categoryID)
		if _, err := C_postCategories.Exec(postID, Category_id); err != nil {
			return fmt.Errorf("error linking post to category %s: %w", categoryID, err)
		}
	}
	return nil
}

func CreateCategorie(name_categorie string) (int, error) {
	var count int

	err := database.DataBase.QueryRow("SELECT COUNT(*) FROM categories WHERE name = ?", name_categorie).Scan(&count)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error checking category existence: %w", err)
	}
	if count > 0 {
		return http.StatusNotFound, errors.New("category already exists")
	}

	C_categories, err := database.DataBase.Prepare(`INSERT INTO categories (name) VALUES (?)`)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error preparing statement: %w", err)
	}
	defer C_categories.Close()

	if _, err := C_categories.Exec(name_categorie); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error executing statement: %w", err)
	}

	return http.StatusOK, nil
}

func CreateComment(comment models.Comment) error {
	C_comment, err := database.DataBase.Prepare(`INSERT INTO comments(post_id, user_id, content) VALUES(?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}

	defer C_comment.Close()

	if _, err = C_comment.Exec(comment.PostId, comment.UserId, comment.Content); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}
	return nil
}

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
