package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"log"
)

/*
GetAllCategories retrieves all rows from the Category table in the database and returns them as a slice of Category structs.
*/
func GetAllCategories() ([]structs.Category, error) {
	rows, err := Client.Query(`SELECT * FROM Category`)
	if err != nil {
		log.Println("Error when selecting all categories: ", err)
		return nil, err
	}

	foundCategories, err2 := categoryRowsToSlice(rows)
	if err2 != nil {
		log.Println("Error when converting rows to slice: ", err2)
		return nil, err2
	}
	return foundCategories, nil
}

/*
GetCategoryByID retrieves a single row from the Category table in the database based on the ID and returns it as a Category struct.
*/
func GetCategoryByID(id string) (*structs.Category, error) {
	exists, err := categoryExists(id)
	if err != nil {
		log.Println("Error when checking if category exists: ", err)
		return nil, err
	}
	if exists {
		var category structs.Category
		err2 := Client.QueryRow("SELECT * FROM Category WHERE ID = ?", id).Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err2 != nil {
			log.Println("Error when selecting category by ID: ", err2)
			return nil, err2
		}
		return &category, nil
	} else {
		log.Println("Category does not exist")
		return nil, fmt.Errorf("category ID does not match any categories in DB")
	}
}

/*
AddCategory adds a single row to the Category table in the database. Returns the ID if successful, or an error if not.
*/
func AddCategory(category structs.Category) (string, error) {
	// Generate and retrieve new UUID
	id, err := GenerateUUID(Client)
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return "", err
	}

	// Insert category
	_, err2 := Client.Exec(
		`INSERT INTO Category (ID, Name, Description) VALUES (?, ?, ?)`,
		id,
		category.Name,
		category.Description,
	)
	if err2 != nil {
		log.Println("Error inserting category: ", err2)
		return "", err2
	}

	return id, nil // Return the UUID
}

/*
UpdateCategory updates a single row in the Category table in the database based on the ID in the provided Category struct.
Returns nil if successful, or an error if not.
*/
func UpdateCategory(category structs.Category) error {
	// Check if category exists
	exists, err := categoryExists(category.ID)
	if err != nil {
		log.Println("Error when checking if category exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec(
			"UPDATE Category SET Name = ?, Description = ? WHERE ID = ?",
			category.Name,
			category.Description,
			category.ID,
		)
		if err2 != nil {
			log.Println("Error when updating category: ", err2)
			return err2
		}
	} else {
		log.Println("Category does not exist")
		return fmt.Errorf("category ID does not match any categories in DB")
	}
	return nil
}

/*
DeleteCategoryByID deletes a single row from the Category table in the database based on the ID.
Returns nil if successful, or an error if not.
*/
func DeleteCategoryByID(id string) error {
	exists, err := categoryExists(id)
	if err != nil {
		log.Println("Error when checking if category exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec("DELETE FROM Category WHERE ID = ?", id)
		if err2 != nil {
			log.Println("Error when deleting category: ", err2)
			return err2
		}
	} else {
		log.Println("Category does not exist")
		return fmt.Errorf("category ID does not match any categories in DB")
	}
	return nil
}

/*
categoryRowsToSlice converts the rows from a SQL query to a slice of Category structs.
*/
func categoryRowsToSlice(rows *sql.Rows) ([]structs.Category, error) {
	var categorySlice []structs.Category
	for rows.Next() {
		var category structs.Category
		err2 := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err2 != nil {
			return nil, err2
		}
		categorySlice = append(categorySlice, category)
	}
	return categorySlice, nil
}

/*
CategoryExists checks if a category with the provided ID exists in the database.
*/
func categoryExists(id string) (bool, error) {
	var exists bool
	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM Category WHERE ID = ?)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
