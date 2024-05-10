package db

import (
	"Database_Project/internal/structs"
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

// AddCartItem adds a new cart item to the database
func AddCartItem(cartItem structs.CartItem) error {
	_, err := Client.Exec("INSERT INTO CartItem (UserAccountID, ProductID, Quantity) VALUES (?, ?, ?)",
		cartItem.UserAccountID, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		log.Println("Error adding cart item:", err)
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
