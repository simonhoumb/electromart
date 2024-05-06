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

// RegisterUser creates a new user in the database.
func (db *UserDB) RegisterUser(userID, username, hashedPassword, email, firstName, lastName string, phone string, cartID string) error {
	query := `
        INSERT INTO User (ID, Username, Password, Email, FirstName, LastName, Phone, CartID)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Client.Exec(query, userID, username, hashedPassword, email, firstName, lastName, phone, cartID)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}
	log.Printf("User successfully inserted into database. UserID: %s", userID)
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

// UpdateUserProfile updates the user's profile in the database, including their address.
func (db *UserDB) UpdateUserProfile(user structs.ActiveUser) error {
	// Begin transaction
	tx, err := db.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			log.Printf("Transaction rolled back: %v", err)
			return
		}
		// Commit the transaction if everything is successful
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
	}()

	// Update User table
	_, err = tx.Exec(`
        UPDATE User SET Email=?, FirstName=?, LastName=?, Phone=? WHERE ID=?
    `, user.Email, user.FirstName, user.LastName, user.Phone, user.ID)
	if err != nil {
		return err
	}

	// Check if an address already exists for the user
	var existingAddressID string
	err = tx.QueryRow("SELECT ID FROM Address WHERE UserID = ? AND UserCartID = ?", user.ID, user.CartID).Scan(&existingAddressID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If the user provided a new address, insert it
	if user.Address.String != "" && user.PostCode.String != "" {
		if existingAddressID != "" {
			// Update existing address
			_, err = tx.Exec(`
				UPDATE Address
				SET Street=?, PostalCode=?
				WHERE UserID=? AND UserCartID=?
			`, user.Address.String, user.PostCode.String, user.ID, user.CartID)
			if err != nil {
				return err
			}
		} else {
			// Insert new address
			_, err = tx.Exec(`
				INSERT INTO Address (ID, Street, PostalCode, UserID, UserCartID)
				VALUES (?, ?, ?, ?, ?)
			`, uuid.New().String(), user.Address.String, user.PostCode.String, user.ID, user.CartID)
			if err != nil {
				return err
			}
		}
	} else {
		// If the user did not provide a new address, delete the existing one if it exists
		if existingAddressID != "" {
			_, err = tx.Exec("DELETE FROM Address WHERE ID = ?", existingAddressID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteUser deletes a user from the database along with their associated address and cart.
func (db *UserDB) DeleteUser(username string) error {
	// Begin transaction
	tx, err := db.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			log.Printf("Transaction rolled back: %v", err)
			return
		}
		// Commit the transaction if everything is successful
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
	}()

	// Get the user's ID
	var userID string
	err = tx.QueryRow("SELECT ID FROM User WHERE Username = ?", username).Scan(&userID)
	if err != nil {
		return err
	}

	// Get the user's CartID
	var cartID string
	err = tx.QueryRow("SELECT CartID FROM User WHERE Username = ?", username).Scan(&cartID)
	if err != nil {
		return err
	}

	// Delete Address associated with the user
	_, err = tx.Exec("DELETE FROM Address WHERE UserID = ?", userID)
	if err != nil {
		return err
	}

	// Delete User
	_, err = tx.Exec("DELETE FROM User WHERE Username = ?", username)
	if err != nil {
		return err
	}

	// Delete Cart associated with the user
	_, err = tx.Exec("DELETE FROM Cart WHERE ID = ?", cartID)
	if err != nil {
		return err
	}

	return nil
}
