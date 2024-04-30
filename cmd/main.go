package main

import (
	server "Database_Project/Internal"
	"Database_Project/db"
	database_2024 "Database_Project/db"
	"Database_Project/utils"
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

// main is the entry point for the program.
func main() {
	database := database_2024.Connect()
	// Close the database connection when the main function returns.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)
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