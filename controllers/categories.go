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
	var categories models.Categories
	if categori == "" {
		rows, err := utils.DataBase.Query("SELECT * FROM categories")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var categoriSlice []models.Categories
		for rows.Next() {
			rows.Scan(&categories.Id, categories.Categori, categories.PostId)
			categoriSlice = append(categoriSlice, categories)
		}
		return categoriSlice, nil
	} else {
		rows, err := utils.DataBase.Query("SELECT categori FROM categories WHERE categori = ?", categori)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var categoriSlice []models.Categories
		for rows.Next() {
			rows.Scan(&categories.Id, categories.Categori, categories.PostId)
			categoriSlice = append(categoriSlice, categories)
		}
		return categoriSlice, nil
	}
}
