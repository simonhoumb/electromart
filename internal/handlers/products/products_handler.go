package products

import (
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

	// In one if statement
	if productsJSON, err := json.Marshal(products); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
		return
	} else {
		w.Header().Set("content-type", "application/json")
		if _, err := w.Write(productsJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
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

func (req CreateProductRequest) Validate() error {
	values := reflect.ValueOf(req)
	for i := 0; i < values.NumField(); i++ {
		switch values.Field(i).Interface().(type) {
		case int:
			if values.Field(i).Int() < 0 {
				return fmt.Errorf("field %s must be a positive integer", values.Type().Field(i).Name)
			}
		default:
			if values.Field(i).IsZero() {
				return fmt.Errorf("field '%s' is invalid and/or required", values.Type().Field(i).Name)
			}
		}
	}
	return nil
}

type CreateProductResponse struct {
	ID string `json:"id"`
}

func handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); utils.HandleError(w, r, http.StatusBadRequest, err, "error during decoding request") {
		return
	}

	if err := req.Validate(); utils.HandleError(w, r, http.StatusBadRequest, err, "invalid request json, check documentation") {
		return
	}

	// Create the product
	// productID := db.somethingsomething()
	productID := "productID"

	// Two above in one if statement
	if productIDJSON, err := json.Marshal(CreateProductResponse{ID: productID}); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
		return
	} else {
		if _, err := w.Write(productIDJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
			return
		}
	}
}
