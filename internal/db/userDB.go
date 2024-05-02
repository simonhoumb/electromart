package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"errors"
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
func (db *UserDB) RegisterUser(userID, username, hashedPassword, email, firstName, lastName string, phone string, cartID string, address sql.NullString, postCode sql.NullString) error {
	query := `INSERT INTO User (ID, Username, Password, Email, FirstName, LastName, Phone, CartID, AddressID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, (SELECT ID FROM Address WHERE Street=? AND PostalCode=?))`
	_, err := db.Client.Exec(query, userID, username, hashedPassword, email, firstName, lastName, phone, cartID, address.String, postCode.String)
	if err != nil {
		return err
	}
	return nil
}

// GetUser retrieves a user with given username from the database.
func (db *UserDB) GetUser(username string) (structs.ActiveUser, error) {
	var user structs.ActiveUser

	query := `SELECT User.ID, User.Username, User.Email, User.Password, User.FirstName, User.LastName, User.Phone, User.CartID,
       	   Address.Street, PostalCode.PostalCode 
              FROM User 
              LEFT JOIN Address ON User.ID = Address.UserID
              LEFT JOIN PostalCode ON Address.PostalCode = PostalCode.PostalCode
              WHERE User.Username = ?`

	err := db.Client.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.CartID,
		&user.Address,
		&user.PostCode,
	)

	if err != nil {
		log.Printf("Error fetching user info for %v: %v", username, err)
		return user, err
	}

	return user, nil
}

// UpdateUserProfile updates the user profile in the database.
func (db *UserDB) UpdateUserProfile(user structs.ActiveUser) error {
	// Validate user input
	if !user.Address.Valid || !user.PostCode.Valid {
		return errors.New("invalid info: both Address and PostCode must be set")
	}

	// Begin transaction
	tx, err := db.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			return
		}
		// Commit the transaction if everything is successful
		err = tx.Commit()
	}()

	// Update User table
	_, err = tx.Exec(`
        UPDATE User SET Email=?, FirstName=?, LastName=? WHERE ID=?
    `, user.Email, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return err
	}

	// Check if the user already has an address
	var existingAddressID string
	err = tx.QueryRow("SELECT ID FROM Address WHERE UserID=? AND UserCartID=?", user.ID, user.CartID).Scan(&existingAddressID)
	if err != nil {
		if err != sql.ErrNoRows {
			// Another error occurred, return it
			return err
		}
		// If no address exists, insert a new address
		addressID := uuid.New().String()

		// Verify the uniqueness of the generated address ID
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM Address WHERE ID=?", addressID).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("generated address ID is not unique")
		}

		// Check if the user exists
		var userCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM User WHERE ID=? AND CartID=?", user.ID, user.CartID).Scan(&userCount)
		if err != nil {
			return err
		}
		if userCount == 0 {
			return errors.New("user does not exist")
		}

		// Insert the new address
		_, err = tx.Exec(`
    INSERT INTO Address (ID, Street, PostalCode, UserID, UserCartID) VALUES (?, ?, ?, ?, ?)
`, addressID, user.Address.String, user.PostCode.String, user.ID, user.CartID)
		if err != nil {
			return err
		}

	} else {
		// If an address exists, update the existing address
		_, err = tx.Exec(`
            UPDATE Address SET Street=?, PostalCode=? WHERE ID=?
        `, user.Address.String, user.PostCode.String, existingAddressID)
		if err != nil {
			return err
		}
	}

	return nil
}
