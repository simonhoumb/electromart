package db

import (
	"Database_Project/internal/session"
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

// GetAllCartItems retrieves all cart items for the currently logged-in user.
func (db *UserDB) GetAllCartItems(r *http.Request) ([]structs.CartItem, error) {
	// 1. Get user ID from session
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		return nil, fmt.Errorf("error getting session: %v", err)
	}

	userIDValue := session.Values["userID"]
	userID, ok := userIDValue.(uuid.UUID)
	if !ok || userID.String() == "" { // Changed from != uuid.Nil to == ""
		return nil, fmt.Errorf("user not logged in or invalid userID")
	}

	// 2. Query the database with userID
	rows, err := db.Client.Query(`
        SELECT UserAccountID, ProductID, Quantity 
        FROM CartItem
        WHERE UserAccountID = ?
    `, userID)
	if err != nil {
		log.Println("Error fetching cart items:", err)
		return nil, err
	}
	defer rows.Close()

	// 3. Process the results
	var cartItems []structs.CartItem
	for rows.Next() {
		var cartItem structs.CartItem
		if err := rows.Scan(&cartItem.UserAccountID, &cartItem.ProductID, &cartItem.Quantity); err != nil {
			log.Println("Error scanning cart item row:", err)
			continue
		}
		cartItems = append(cartItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over cart item rows:", err)
		return nil, err
	}

	return cartItems, nil
}

// GetCartItemsByUser retrieves all cart items for a specific user from the database
func GetCartItemsByUser(userID string) ([]structs.CartItem, error) {
	rows, err := Client.Query(`SELECT UserAccountID, ProductID, Quantity FROM CartItem WHERE UserAccountID = ?`, userID)
	if err != nil {
		log.Println("Error fetching cart items by user:", err)
		return nil, err
	}
	defer rows.Close()

	var cartItems []structs.CartItem
	for rows.Next() {
		var cartItem structs.CartItem
		if err := rows.Scan(&cartItem.UserAccountID, &cartItem.ProductID, &cartItem.Quantity); err != nil {
			log.Println("Error scanning cart item row:", err)
			continue
		}
		cartItems = append(cartItems, cartItem)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over cart item rows:", err)
		return nil, err
	}

	return cartItems, nil
}

func AddCartItem(cartItem structs.CartItem) error {
	// Check if the item already exists for the user
	existingItem, err := GetCartItemByUserIDAndProductID(cartItem.UserAccountID, cartItem.ProductID)
	if err != nil && err != sql.ErrNoRows { // Error other than no rows found
		log.Println("Error checking for existing cart item:", err)
		return err
	}

	if existingItem != nil {
		// Update quantity if the item exists
		newQuantity := existingItem.Quantity + cartItem.Quantity
		_, err = Client.Exec("UPDATE CartItem SET Quantity = ? WHERE UserAccountID = ? AND ProductID = ?",
			newQuantity, cartItem.UserAccountID, cartItem.ProductID)
	} else {
		// Insert a new cart item if it doesn't exist
		_, err = Client.Exec("INSERT INTO CartItem (UserAccountID, ProductID, Quantity) VALUES (?, ?, ?)",
			cartItem.UserAccountID, cartItem.ProductID, cartItem.Quantity)
	}

	if err != nil {
		log.Println("Error adding/updating cart item:", err)
		return err
	}

	return nil
}

// GetCartItemByUserIDAndProductID retrieves a cart item by UserAccountID and ProductID
func GetCartItemByUserIDAndProductID(userID, productID string) (*structs.CartItem, error) {
	var cartItem structs.CartItem
	err := Client.QueryRow("SELECT * FROM CartItem WHERE UserAccountID = ? AND ProductID = ?", userID, productID).Scan(
		&cartItem.UserAccountID, &cartItem.ProductID, &cartItem.Quantity,
	)
	if err != nil {
		// If no rows were found, return nil, sql.ErrNoRows to differentiate from other errors
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Println("Error retrieving cart item:", err)
		return nil, err
	}
	return &cartItem, nil
}

func UpdateCartItemQuantity(userID, productID string, newQuantity int) error {
	_, err := Client.Exec("UPDATE CartItem SET Quantity = ? WHERE UserAccountID = ? AND ProductID = ?", newQuantity, userID, productID)
	if err != nil {
		log.Println("Error updating cart item quantity:", err)
		return err
	}
	return nil
}

// DeleteCartItem removes a cart item from the database
func DeleteCartItem(userID, productID string) error {
	_, err := Client.Exec("DELETE FROM CartItem WHERE UserAccountID = ? AND ProductID = ?", userID, productID)
	if err != nil {
		log.Println("Error deleting cart item:", err)
		return err
	}
	return nil
}
