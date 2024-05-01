package products

import (
	"Database_Project/internal/db"
	"Database_Project/internal/structs"
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// Implemented methods for the endpoint
var productsImplementedMethods = []string{
	http.MethodGet,
	http.MethodPost,
}

/*
HandleProducts for the /products endpoint.
*/
func HandleProducts(w http.ResponseWriter, r *http.Request) {
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
				productsImplementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleGetAllRequest(w http.ResponseWriter, r *http.Request) {
	// Get all products
	products, err := db.GetAllProducts(db.Client)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error getting products from database") {
		return
	}

	// Return the products
	if productsJSON, err := json.Marshal(products); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
		return
	} else {
		if _, err := w.Write(productsJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
			return
		}
	}
}

func handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var product structs.Product

	if err := json.NewDecoder(r.Body).Decode(&product); utils.HandleError(w, r, http.StatusBadRequest, err, "error during decoding request") {
		return
	}

	if err := product.ValidateNewProductRequest(); utils.HandleError(w, r, http.StatusBadRequest, err, "invalid request json, check documentation") {
		return
	}

	// Create the product
	productID, err := db.AddProduct(db.Client, product)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error adding product to database") {
		return
	}

	// Two above in one if statement
	if productIDJSON, err := json.Marshal(structs.CreateProductResponse{ID: productID}); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
		return
	} else {
		if _, err := w.Write(productIDJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
			return
		}
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
