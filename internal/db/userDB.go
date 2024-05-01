package db

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserDB struct {
	Client *sql.DB
}

// UserExists checks if a user with given username and password exists in the database.
func (db *UserDB) UserExists(username string, password string) (bool, error) {
	// Query to check if the user exists
	queryStmt := `SELECT COUNT(1) FROM User WHERE username=?`

	// Use QueryRow to return a row and scan the returned id into the User struct
	var exists bool
	err := db.Client.QueryRow(queryStmt, username, password).Scan(&exists)

	if err != nil {
		// If an error is returned from the query, return a more contextual error message and the error
		return false, fmt.Errorf("unable to check if user exists: %v", err)
	}

	return exists, nil
}

// CheckLogin checks if the given username and password match a record in the database.
func (db *UserDB) CheckLogin(username string, password string) (bool, error) {
	queryStmt := `SELECT password FROM User WHERE username=?`

	var hashedPassword string
	err := db.Client.QueryRow(queryStmt, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username was not found
			return false, fmt.Errorf("username not found")
		}
		// Another error occurred
		return false, fmt.Errorf("internal error when fetching row: %v", err)
	}

	// this will check if the password hashes match, returns an error if they don't
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Password does not match stored hash
		return false, fmt.Errorf("unauthorized")
	}

	// Both username and password match the record in the database
	return true, nil
}
