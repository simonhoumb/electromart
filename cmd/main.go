package main

import (
	database_2024 "Database_Project"
	server "Database_Project/Internal"
	"Database_Project/db"
	"Database_Project/structs"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

// main is the entry point for the program.
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
