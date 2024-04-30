package products

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
TODO:
- Create a new product
- Update a product
- Delete a product
- Get all products
    - Filter by:
        - Category
		- Price range
		- Name
		- Brand
- Get a single product (by ID)
*/

// Implemented methods for the endpoint
var implementedMethods = []string{
	http.MethodGet,
}

/*
Handler for the /products endpoint.
*/
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleGetAllRequest(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleGetAllRequest(w http.ResponseWriter, r *http.Request) {
	// Get all products
	// products := db.somethingsomething()
	products := []string{"product1", "product2", "product3"}

	// Marshal the products into a JSON object
	productsJSON, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON object to the response
	_, err = w.Write(productsJSON)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}
