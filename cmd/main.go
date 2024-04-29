package main

import (
	"Database_Project/db"
	database_2024 "Database_Project/db"
	"Database_Project/structs"
	"Database_Project/utils"
	"database/sql"
	"encoding/json"
	"net/http"
)

var database *sql.DB

func main() {
	database = database_2024.Connect()
	defer database.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/login", utils.CheckLogin(database))
	http.HandleFunc("/logout", utils.LogoutUser(database))
	http.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		utils.GetCartItems(w, r, database)
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		utils.GetUserProfile(w, r, database)
	})
	http.HandleFunc("/api/categories", db.GetCategoriesHandler(database))
	http.HandleFunc("/register", utils.RegisterUser(database))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/loginPage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/Login.html")
	})

	http.HandleFunc("/registerPage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/register.html")
	})

	http.ListenAndServe(":8080", nil)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM Product")
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
