package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

var Client *sql.DB

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "dev:Devdevdev2!@tcp(10.212.169.49:3306)/ElectroMart")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database.")
	return db
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetCategoriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categories []Category

		rows, err := db.Query("SELECT id, name, description FROM Category")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var c Category
			if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			categories = append(categories, c)
		}

		jsonData, err := json.Marshal(categories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

// getDataSourceName returns a Data Source Name string for connecting to a MySQL database.
func getDataSourceName() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbNet := os.Getenv("DB_NET")
	dbIP := os.Getenv("DB_IP")
	dbName := os.Getenv("DB_NAME")

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  dbNet,
		Addr:                 dbIP,
		DBName:               dbName,
		AllowNativePasswords: true,
	}

	return cfg.FormatDSN()
}

// openDatabaseConnection opens a connection to the database and returns a pointer to the database.
func openDatabaseConnection() *sql.DB {
	// Open a database connection.
	db, err := sql.Open("mysql", getDataSourceName())
	if err != nil {
		log.Fatal("Error when opening the database: ", err)
	}
	return db
}

// pingDatabase pings the database to verify the connection.
func pingDatabase(db *sql.DB) {
	// Ping the database to verify the connection.
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error when pinging database: ", err)
	}
}
