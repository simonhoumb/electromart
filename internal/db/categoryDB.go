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
	rows, err := Client.Query(`SELECT * FROM ElectroMartDB.Category`)
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
GetCategoryByName retrieves a single row from the Category table in the database based on the name and returns it as a Category struct.
*/
func GetCategoryByName(name string) (*structs.Category, error) {
	exists, err := categoryExists(name)
	if err != nil {
		log.Println("Error when checking if category exists: ", err)
		return nil, err
	}
	if exists {
		var category structs.Category
		err2 := Client.QueryRow("SELECT * FROM ElectroMartDB.Category WHERE Name = ?", name).Scan(
			&category.Name,
			&category.Description,
		)
		if err2 != nil {
			log.Println("Error when selecting category by name: ", err2)
			return nil, err2
		}
		return &category, nil
	} else {
		log.Println("Category does not exist")
		return nil, fmt.Errorf("category name does not match any categories in DB")
	}
}

/*
AddCategory adds a single row to the Category table in the database. Returns the name if successful, or an error if not.
*/
func AddCategory(category structs.Category) error {
	// Insert category
	_, err2 := Client.Exec(
		`INSERT INTO ElectroMartDB.Category (Name, Description) VALUES (?, ?)`,
		category.Name,
		category.Description,
	)
	if err2 != nil {
		log.Println("Error inserting category: ", err2)
		return err2
	}

	return nil
}

/*
UpdateCategory updates a single row in the Category table in the database based on the name in the provided Category struct.
Returns nil if successful, or an error if not.
*/
func UpdateCategory(category structs.Category) error {
	// Check if category exists
	exists, err := categoryExists(category.Name)
	if err != nil {
		log.Println("Error when checking if category exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec(
			"UPDATE ElectroMartDB.Category SET Description = ? WHERE Name = ?",
			category.Description,
			category.Name,
		)
		if err2 != nil {
			log.Println("Error when updating category: ", err2)
			return err2
		}
	} else {
		log.Println("Category does not exist")
		return fmt.Errorf("category name does not match any categories in DB")
	}
	return nil
}

/*
DeleteCategoryByName deletes a single row from the Category table in the database based on the name.
Returns nil if successful, or an error if not.
*/
func DeleteCategoryByName(name string) error {
	exists, err := categoryExists(name)
	if err != nil {
		log.Println("Error when checking if category exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec("DELETE FROM ElectroMartDB.Category WHERE Name = ?", name)
		if err2 != nil {
			log.Println("Error when deleting category: ", err2)
			return err2
		}
	} else {
		log.Println("Category does not exist")
		return fmt.Errorf("category name does not match any categories in DB")
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
CategoryExists checks if a category with the provided name exists in the database.
*/
func categoryExists(name string) (bool, error) {
	var exists bool
	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM ElectroMartDB.Category WHERE Name = ?)`, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
