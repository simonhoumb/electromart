package main

import (
	database_2024 "Database_Project"
	"database/sql"
	"log"
)

func main() {
	database := database_2024.Connect()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)

}
