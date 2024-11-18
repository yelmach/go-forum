package utils

import (
	"database/sql"
	"fmt"

	"forum/database"
)

func VerifyCategoriesMatch(categories []string) error {
	dbCategories, err := database.DataBase.Query(`SELECT name FROM categories`)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}
	defer dbCategories.Close()

	categoriesFromDb := make(map[string]bool)

	for dbCategories.Next() {
		var category string
		if err := dbCategories.Scan(&category); err != nil {
			return err
		}
		categoriesFromDb[category] = true
	}

	for _, category := range categories {
		if !categoriesFromDb[category] {
			return fmt.Errorf("category '%s' not found in the database", category)
		}
	}
	return nil
}

func HasUniqueCategories(categories []string) bool {
	isDouble := make(map[string]bool)
	for _, category := range categories {
		if isDouble[category] {
			return false
		}
		isDouble[category] = true
	}
	return true
}