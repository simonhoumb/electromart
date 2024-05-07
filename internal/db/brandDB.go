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
GetBrandByName retrieves a single row from the Brand table in the database based on the ID and returns it as a Brand struct.
*/
func GetBrandByName(name string) (*structs.Brand, error) {
	exists, err := BrandExists(name)
	if err != nil {
		log.Println("Error when checking if Brand exists: ", err)
		return nil, err
	}
	if exists {
		var Brand structs.Brand
		err2 := Client.QueryRow("SELECT * FROM Brand WHERE Name = ?", name).Scan(
			&Brand.Name,
			&Brand.Description,
		)
		if err2 != nil {
			log.Println("Error when selecting Brand by name: ", err2)
			return nil, err2
		}
		return &Brand, nil
	} else {
		log.Println("Brand does not exist")
		return nil, fmt.Errorf("Brand name does not match any Brands in DB")
	}
}

/*
AddBrand adds a single row to the Brand table in the database. Returns the ID if successful, or an error if not.
*/
func AddBrand(Brand structs.Brand) error {
	// Insert Brand
	_, err2 := Client.Exec(
		`INSERT INTO Brand (Name, Description) VALUES (?, ?)`,
		Brand.Name,
		Brand.Description,
	)
	if err2 != nil {
		log.Println("Error inserting Brand: ", err2)
		return err2
	}

	return nil
}

/*
UpdateBrand updates a single row in the Brand table in the database based on the ID in the provided Brand struct.
Returns nil if successful, or an error if not.
*/
func UpdateBrand(Brand structs.Brand) error {
	// Check if Brand exists
	exists, err := BrandExists(Brand.Name)
	if err != nil {
		log.Println("Error when checking if Brand exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec(
			"UPDATE Brand SET Description = ? WHERE Name = ?",
			Brand.Description,
			Brand.Name,
		)
		if err2 != nil {
			log.Println("Error when updating Brand: ", err2)
			return err2
		}
	} else {
		log.Println("Brand does not exist")
		return fmt.Errorf("Brand name does not match any Brands in DB")
	}
	return nil
}

/*
DeleteBrandByName deletes a single row from the Brand table in the database based on the ID.
Returns nil if successful, or an error if not.
*/
// TODO: Fix this:
/*
Error when deleting Brand:  Error 1451 (23000): Cannot delete or update a parent row: a foreign key constraint fails (`ElectroMart`.`Product`, CONSTRAINT `FKProduct960541` FOREIGN KEY (`BrandName`) REFERENCES `Brand` (`ID`))
*/
func DeleteBrandByName(name string) error {
	exists, err := BrandExists(name)
	if err != nil {
		log.Println("Error when checking if Brand exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec("DELETE FROM Brand WHERE Name = ?", name)
		if err2 != nil {
			log.Println("Error when deleting Brand: ", err2)
			return err2
		}
	} else {
		log.Println("Brand does not exist")
		return fmt.Errorf("Brand name does not match any Brands in DB")
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
func BrandExists(name string) (bool, error) {
	var exists bool
	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM Brand WHERE Name = ?)`, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
