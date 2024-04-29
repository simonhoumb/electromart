package main

import (
	database_2024 "Database_Project"
	"database/sql"
	"log"
)

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
