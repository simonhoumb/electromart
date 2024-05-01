package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"log"
	"strconv"
	"strings"
)

/*
SearchProducts retrieves rows from the Product table in the database based on the query string provided.
The query string is used to search for products by name, description, brand name, and category name.
Returns a slice of Product structs if successful, or an error if not.
*/
func SearchProducts(db *sql.DB, query string) ([]structs.Product, error) {
	lowerQuery := strings.ToLower(query) // Convert query to lowercase
	rows, err := db.Query(
		`
      SELECT * FROM Product
      WHERE LOWER(Name) LIKE ? OR Description LIKE ? OR BrandID IN (SELECT ID FROM Brand WHERE LOWER(Name) LIKE ?)
        AND CategoryID IN (SELECT ID FROM Category WHERE LOWER(Name) LIKE ?);
  `, "%"+lowerQuery+"%", "%"+lowerQuery+"%", "%"+lowerQuery+"%", "%"+lowerQuery+"%",
	)
	if err != nil {
		log.Println("Error when querying for products: ", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err2 := rows.Close()
		if err2 != nil {
			log.Println("Error when closing rows: ", err2)
		}
	}(rows)

	products, err3 := rowsToSlice(rows)
	if err3 != nil {
		log.Println("Error when converting rows to slice: ", err3)
		return nil, err3
	}
	return products, nil
}

/*
GetAllProducts retrieves all rows from the Product table in the database and returns them as a slice of Product structs.
*/
func GetAllProducts(db *sql.DB) ([]structs.Product, error) {
	rows, err := db.Query(`SELECT * FROM Product`)
	if err != nil {
		log.Println("Error when selecting all products: ", err)
		return nil, err
	}

	foundProducts, err2 := rowsToSlice(rows)
	if err2 != nil {
		log.Println("Error when converting rows to slice: ", err2)
		return nil, err2
	}
	return foundProducts, nil
}

/*
GetProductByID retrieves a single row from the Product table in the database based on the ID and returns it as a Product struct.
*/
func GetProductByID(db *sql.DB, id string) (structs.Product, error) {
	var product structs.Product
	err := db.QueryRow("SELECT * FROM Product WHERE ID = ?", id).Scan(
		&product.ID,
		&product.Name,
		&product.BrandID,
		&product.CategoryID,
		&product.Description,
		&product.QtyInStock,
		&product.Price,
	)
	if err != nil {
		log.Println("Error when selecting product by ID: ", err)
		return product, err
	}
	return product, nil
}

/*
AddProduct adds a single row to the Product table in the database. Returns the ID if successful.
*/
func AddProduct(db *sql.DB, product structs.Product) (string, error) {
	result, err := db.Exec(
		`INSERT INTO Product (ID, Name, BrandID, CategoryID, Description, QtyInStock, Price) VALUES (UUID(), ?, ?, ?,
?, ?, ?)`,
		product.Name,
		product.BrandID,
		product.CategoryID,
		product.Description,
		product.QtyInStock,
		product.Price,
	)
	if err != nil {
		log.Println("Error when adding product: ", err)
		return "", err
	}
	id, err2 := result.LastInsertId()
	if err2 != nil {
		log.Println("Error when getting last insert ID: ", err2)
		return "", err2
	}

	return strconv.FormatInt(id, 10), nil
}

/*
UpdateProduct updates a single row in the Product table in the database based on the ID in the provided Product struct.
Returns nil if successful, or an error if not.
*/
func UpdateProduct(db *sql.DB, product structs.Product) error {
	_, err := db.Exec(
		"UPDATE Product SET Name = ?, BrandID = ?, CategoryID = ?, Description = ?, QtyInStock = ?, Price = ? WHERE ID = ?",
		product.Name,
		product.BrandID,
		product.CategoryID,
		product.Description,
		product.QtyInStock,
		product.Price,
		product.ID,
	)
	if err != nil {
		log.Println("Error when updating product: ", err)
		return err
	}
	return nil
}

/*
DeleteProduct deletes a single row from the Product table in the database based on the ID. Returns nil if successful, or an error if not.
*/
func DeleteProductByID(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM Product WHERE ID = ?", id)
	if err != nil {
		log.Println("Error when deleting product: ", err)
		return err
	}
	return nil
}

func rowsToSlice(rows *sql.Rows) ([]structs.Product, error) {
	var productSlice []structs.Product
	for rows.Next() {
		var product structs.Product
		err2 := rows.Scan(
			&product.ID,
			&product.Name,
			&product.BrandID,
			&product.CategoryID,
			&product.Description,
			&product.QtyInStock,
			&product.Price,
		)
		if err2 != nil {
			return nil, err2
		}
		productSlice = append(productSlice, product)
	}
	return productSlice, nil
}

func ProductExists(product structs.Product) (bool, error) {
	var exists bool

	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM Product WHERE ID = ?)`, product.ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
