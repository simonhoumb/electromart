package main

import (
	server "Database_Project/internal"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var Store *sessions.CookieStore

// init loads the .env file when the program starts.
func init() {
	Store = sessions.NewCookieStore([]byte("your-very-secret-key"))
	Store.Options = &sessions.Options{
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

// main is the entry point for the program.
func main() {
	server.Start()
}
