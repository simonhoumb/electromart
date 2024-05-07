package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"log"
)

/*
GetAllBrands retrieves all rows from the Brand table in the database and returns them as a slice of Brand structs.
*/
func GetAllBrands() ([]structs.Brand, error) {
	rows, err := Client.Query(`SELECT * FROM Brand`)
	if err != nil {
		log.Println("Error when selecting all Brands: ", err)
		return nil, err
	}

	foundBrands, err2 := BrandRowsToSlice(rows)
	if err2 != nil {
		log.Println("Error when converting rows to slice: ", err2)
		return nil, err2
	}
	return foundBrands, nil
}

/*
GetBrandByID retrieves a single row from the Brand table in the database based on the ID and returns it as a Brand struct.
*/
func GetBrandByID(id string) (*structs.Brand, error) {
	exists, err := BrandExists(id)
	if err != nil {
		log.Println("Error when checking if Brand exists: ", err)
		return nil, err
	}
	if exists {
		var Brand structs.Brand
		err2 := Client.QueryRow("SELECT * FROM Brand WHERE ID = ?", id).Scan(
			&Brand.ID,
			&Brand.Name,
			&Brand.Description,
		)
		if err2 != nil {
			log.Println("Error when selecting Brand by ID: ", err2)
			return nil, err2
		}
		return &Brand, nil
	} else {
		log.Println("Brand does not exist")
		return nil, fmt.Errorf("Brand ID does not match any Brands in DB")
	}
}

/*
AddBrand adds a single row to the Brand table in the database. Returns the ID if successful, or an error if not.
*/
func AddBrand(Brand structs.Brand) (string, error) {
	// Generate and retrieve new UUID
	id, err := GenerateUUID(Client)
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return "", err
	}

	// Insert Brand
	_, err2 := Client.Exec(
		`INSERT INTO Brand (ID, Name, Description) VALUES (?, ?, ?)`,
		id,
		Brand.Name,
		Brand.Description,
	)
	if err2 != nil {
		log.Println("Error inserting Brand: ", err2)
		return "", err2
	}

	return id, nil // Return the UUID
}

/*
UpdateBrand updates a single row in the Brand table in the database based on the ID in the provided Brand struct.
Returns nil if successful, or an error if not.
*/
func UpdateBrand(Brand structs.Brand) error {
	// Check if Brand exists
	exists, err := BrandExists(Brand.ID)
	if err != nil {
		log.Println("Error when checking if Brand exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec(
			"UPDATE Brand SET Name = ?, Description = ? WHERE ID = ?",
			Brand.Name,
			Brand.Description,
			Brand.ID,
		)
		if err2 != nil {
			log.Println("Error when updating Brand: ", err2)
			return err2
		}
	} else {
		log.Println("Brand does not exist")
		return fmt.Errorf("Brand ID does not match any Brands in DB")
	}
	return nil
}

/*
DeleteBrandByID deletes a single row from the Brand table in the database based on the ID.
Returns nil if successful, or an error if not.
*/
// TODO: Fix this:
/*
Error when deleting Brand:  Error 1451 (23000): Cannot delete or update a parent row: a foreign key constraint fails (`ElectroMart`.`Product`, CONSTRAINT `FKProduct960541` FOREIGN KEY (`BrandName`) REFERENCES `Brand` (`ID`))
*/
func DeleteBrandByID(id string) error {
	exists, err := BrandExists(id)
	if err != nil {
		log.Println("Error when checking if Brand exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec("DELETE FROM Brand WHERE ID = ?", id)
		if err2 != nil {
			log.Println("Error when deleting Brand: ", err2)
			return err2
		}
	} else {
		log.Println("Brand does not exist")
		return fmt.Errorf("Brand ID does not match any Brands in DB")
	}
	return nil
}

/*
BrandRowsToSlice converts the rows from a SQL query to a slice of Brand structs.
*/
func BrandRowsToSlice(rows *sql.Rows) ([]structs.Brand, error) {
	var BrandSlice []structs.Brand
	for rows.Next() {
		var Brand structs.Brand
		err2 := rows.Scan(
			&Brand.ID,
			&Brand.Name,
			&Brand.Description,
		)
		if err2 != nil {
			return nil, err2
		}
		BrandSlice = append(BrandSlice, Brand)
	}
	return BrandSlice, nil
}

/*
BrandExists checks if a Brand with the provided ID exists in the database.
*/
func BrandExists(id string) (bool, error) {
	var exists bool
	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM Brand WHERE ID = ?)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
