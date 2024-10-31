package controllers

import (
	"errors"

	"forum/models"
	"forum/utils"
)

func CreateCategories(categories models.Categories) (bool, error) {
	var count int
	err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM categories WHERE categori = ? nameCategori", categories.Categori).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, errors.New("categori already exist")
	}
	stmt, err := utils.DataBase.Prepare("INSERT INTO categories (categori,post_id) VALUES(?,?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(categories.Categori, categories.PostId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DisplayCategories(categori string) ([]models.Categories, error) {
	var categories []models.Categories

	var query string
	var args []interface{}

	if categori == "" {
		query = "SELECT * FROM categories"
	} else {
		query = "SELECT * FROM categories WHERE categori = ?"
		args = append(args, categori)
	}

	rows, err := utils.DataBase.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Categories
		err := rows.Scan(&category.Id, &category.Categori, &category.PostId)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
