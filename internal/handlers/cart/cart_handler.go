package cart

import (
	"Database_Project/internal/db"
	"Database_Project/internal/session"
	"Database_Project/internal/structs"
	"encoding/json"
	"fmt"
	"net/http"
)

// HandleCart for the /cart endpoint.
func HandleCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if the user is logged in
	session.CheckSession(func(w http.ResponseWriter, r *http.Request) {
		// Your logic for handling the cart request goes here
		handleCartRequest(w, r)
	})(w, r)
}

// Implemented methods for the endpoint
var implementedMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodDelete,
}

func handleCartRequest(w http.ResponseWriter, r *http.Request) {
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, r)
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodPatch:
		handlePatchRequest(w, r)
	case http.MethodDelete:
		handleDeleteRequest(w, r)
	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	// 1. Check the user session
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	userIDValue := session.Values["userID"].(string)

	// 2. Get cart items for the specific user using userID
	cartItems, err := db.GetCartItemsByUser(userIDValue)
	if err != nil {
		http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
		return
	}

	// 3. Return the cart items
	if err := json.NewEncoder(w).Encode(cartItems); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the session cookie
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	// Log session values for debugging
	fmt.Printf("Session Values: %+v\n", session.Values)

	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Failed to retrieve user ID from session", http.StatusInternalServerError)
		return
	}

	// Decode the cart item from the request body
	var cartItem structs.CartItem
	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	// Set the user ID for the cart item
	cartItem.UserAccountID = userID

	// Print the decoded cartItem struct to the console
	fmt.Printf("Decoded cartItem: %+v\n", cartItem)

	// Add the cart item
	if err := db.AddCartItem(cartItem); err != nil {
		http.Error(w, "Failed to add cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlePatchRequest(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Failed to retrieve user ID from session", http.StatusInternalServerError)
		return
	}

	productID := r.URL.Query().Get("productID")
	if productID == "" {
		http.Error(w, "Missing productID query parameter", http.StatusBadRequest)
		return
	}

	// Decode the update request
	var updateRequest struct {
		NewQuantity int `json:"newQuantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	// Validate newQuantity (e.g., ensure it's not negative)
	if updateRequest.NewQuantity < 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Update the cart item quantity in the database
	if err := db.UpdateCartItemQuantity(userID, productID, updateRequest.NewQuantity); err != nil {
		http.Error(w, "Failed to update cart item quantity", http.StatusInternalServerError)
		return
	}

	// Fetch and return updated cart items to the client
	cartItems, err := db.GetCartItemsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to fetch updated cart items", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(cartItems); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the session cookie
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Failed to retrieve user ID from session", http.StatusInternalServerError)
		return
	}

	productID := r.URL.Query().Get("productID")

	if productID == "" {
		http.Error(w, "Missing productID query parameter", http.StatusBadRequest)
		return
	}

	// Delete the cart item
	if err := db.DeleteCartItem(userID, productID); err != nil {
		http.Error(w, "Failed to delete cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
