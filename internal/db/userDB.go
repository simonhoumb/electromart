package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func (db *UserDB) CreateUserCart() (string, error) {
	cartID := uuid.New().String()
	query := `INSERT INTO Cart (ID, TotalAmount) VALUES (?, ?)`
	_, err := db.Client.Exec(query, cartID, 0)
	if err != nil {
		return "", err
	}
	return cartID, nil
}

// CreateUser creates a new user in the database.
func (db *UserDB) RegisterUser(userID, username, hashedPassword, email, firstName, lastName string, phone string, cartID string) error {
	query := `INSERT INTO User (ID, Username, Password, Email, FirstName, LastName, Phone, CartID) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Client.Exec(query, userID, username, hashedPassword, email, firstName, lastName, phone, cartID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser retrieves a user with given username from the database.
func (db *UserDB) GetUser(username string) (structs.ActiveUser, error) {
	var user structs.ActiveUser

	query := `SELECT User.ID, User.Username, User.Email, User.Password, User.FirstName, User.LastName, 
       	   Address.Street, PostalCode.PostalCode
              FROM User 
              LEFT JOIN Address ON User.ID = Address.UserID
              LEFT JOIN PostalCode ON Address.PostalCode = PostalCode.PostalCode
              WHERE User.Username = ?`

	err := db.Client.QueryRow(query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Address,
		&user.PostCode,
	)

	if err != nil {
		log.Printf("Error fetching user info for %v: %v", username, err)
		return user, err
	}

	return user, nil
}

// UpdateUser updates the user's information in the database.
func (db *UserDB) UpdateUser(user structs.ActiveUser) error {
	// Query to update a user in the database
	queryStmt := `UPDATE User SET FirstName = ?, LastName = ?, Address = ?, PostCode = ? WHERE Username = ?`

	// Execute the SQL command to update the user
	_, err := db.Client.Exec(queryStmt, &user.FirstName, &user.LastName, &user.Address, &user.PostCode, &user.Username)
	if err != nil {
		// Return an error with more context.
		return fmt.Errorf("Failed to update user: %v", err)
	}

	return nil
}
