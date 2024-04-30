package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

var Client *sql.DB

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

// OpenDatabaseConnection opens a connection to the database and returns a pointer to the database.
func OpenDatabaseConnection() *sql.DB {
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
