package main

import (
	server "Database_Project/internal"
	"github.com/joho/godotenv"
	"log"
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
