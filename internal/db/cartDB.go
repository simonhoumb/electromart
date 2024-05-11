package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"log"
)

// GetAllCartItems retrieves all cart items from the database and returns them as a slice of CartItem structs.
func GetAllCartItems() ([]structs.CartItem, error) {
	rows, err := Client.Query(`SELECT UserAccountID, ProductID, Quantity FROM CartItem`)
	if err != nil {
		log.Println("Error fetching all cart items:", err)
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
