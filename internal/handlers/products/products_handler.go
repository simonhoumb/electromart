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
	http.MethodPost,
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

	case http.MethodPost:
		handleCreateRequest(w, r)

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

// TODO: move these structs to a separate file
type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BrandID     string `json:"brand_id"`
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
	QtyInStock  int    `json:"qty_in_stock"`
	Price       int    `json:"price"`
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	BrandID     string `json:"brand_id"`
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
	QtyInStock  int    `json:"qty_in_stock"`
	Price       int    `json:"price"`
}

type CreateProductResponse struct {
	ID string `json:"id"`
}

func handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Create the product
	// productID := db.somethingsomething()
	productID := "productID"

	// Marshal the product ID into a JSON object
	productIDJSON, err := json.Marshal(CreateProductResponse{ID: productID})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON object to the response
	_, err = w.Write(productIDJSON)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}
