package database_2024

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Connect connects to the database and returns a pointer to the database.
func Connect() *sql.DB {
	// Open a database connection.
	db := openDatabaseConnection()

	// Set the maximum number of open and idle connections and the maximum connection lifetime.
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Ping the database to verify the connection.
	pingDatabase(db)

	// Print a message indicating that the connection was successful.
	fmt.Println("Successfully connected to the database.")
	return db
}

// getDataSourceName returns a Data Source Name string for connecting to a MySQL database.
func getDataSourceName() string {
	// Get the database credentials from the environment variables.
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbNet := os.Getenv("DB_NET")
	dbIP := os.Getenv("DB_IP")
	dbName := os.Getenv("DB_NAME")

	// Create a Data Source Name string.
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
