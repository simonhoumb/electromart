package main

import (
	server "Database_Project/Internal"
	"Database_Project/db"
	"Database_Project/structs"
	"encoding/json"
	"net/http"
)

func main() {
	server.Start()
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Client.Query("SELECT * FROM Product")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := make([]structs.Product, 0)
	for rows.Next() {
		product := structs.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	jsonData, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
