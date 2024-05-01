package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"log"
)

func GetAllProducts(db *sql.DB) ([]structs.Product, error) {
	var foundProducts []structs.Product
	rows, err := db.Query("SELECT * FROM Product")
	if err != nil {
		log.Println("Error when selecting all products: ", err)
		return nil, err
	}
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
			log.Fatal("Error when scanning products: ", err2)
			return nil, err2
		}
		foundProducts = append(foundProducts, product)
	}
	return foundProducts, nil
}
