package main

import (
	database_2024 "Database_Project"
	"database/sql"
	"log"

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
	// Connect to the database.
	database := database_2024.Connect()
	// Close the database connection when the main function returns.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)
}
