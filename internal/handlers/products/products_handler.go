package products

import (
	"Database_Project/internal/db"
	"Database_Project/internal/structs"
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
HandleProducts for the /products endpoint.
*/
func HandleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	category := r.PathValue("category")
	brand := r.PathValue("brand")

	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		if category != "" {
			handleGetAllByCategoryRequest(w, r) // Handle by category
		} else if brand != "" {
			handleGetAllByBrandRequest(w, r) // Handle by brand
		} else {
			handleGetAllRequest(w, r) // Default: handle all products
		}
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
	products, err := db.GetAllProducts()
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error getting products from database") {
		return
	}

	// Return the products
	if productsJSON, err := json.Marshal(products); utils.HandleError(
		w,
		r,
		http.StatusInternalServerError,
		err,
		"error during encoding response",
	) {
		return
	} else {
		if _, err := w.Write(productsJSON); utils.HandleError(
			w,
			r,
			http.StatusInternalServerError,
			err,
			"error writing response",
		) {
			return
		}
	}
}

func handleGetAllByCategoryRequest(w http.ResponseWriter, r *http.Request) {
	category, err := utils.GetCategoryFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting category from request") {
		return
	}

	// Get all products by category
	products, err := db.GetAllProductsByCategory(category)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error getting products from database") {
		return
	}

	// Return the products
	if productsJSON, err := json.Marshal(products); utils.HandleError(
		w,
		r,
		http.StatusInternalServerError,
		err,
		"error during encoding response",
	) {
		return
	} else {
		if _, err := w.Write(productsJSON); utils.HandleError(
			w,
			r,
			http.StatusInternalServerError,
			err,
			"error writing response",
		) {
			return
		}
	}
}

func handleGetAllByBrandRequest(w http.ResponseWriter, r *http.Request) {
	brand, err := utils.GetBrandFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting brand from request") {
		return
	}

	// Get all products by category
	products, err := db.GetAllProductsByBrand(brand)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error getting products from database") {
		return
	}

	// Return the products
	if productsJSON, err := json.Marshal(products); utils.HandleError(
		w,
		r,
		http.StatusInternalServerError,
		err,
		"error during encoding response",
	) {
		return
	} else {
		if _, err := w.Write(productsJSON); utils.HandleError(
			w,
			r,
			http.StatusInternalServerError,
			err,
			"error writing response",
		) {
			return
		}
	}
}

func handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var product structs.Product

	if err := json.NewDecoder(r.Body).Decode(&product); utils.HandleError(
		w,
		r,
		http.StatusBadRequest,
		err,
		"error during decoding request",
	) {
		return
	}

	if err := product.ValidateNewProductRequest(); utils.HandleError(
		w,
		r,
		http.StatusBadRequest,
		err,
		"invalid request json, check documentation",
	) {
		return
	}

	// Create the product
	productID, err := db.AddProduct(product)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error adding product to database") {
		return
	}

	// Two above in one if statement
	if productIDJSON, err := json.Marshal(structs.CreateProductResponse{ID: productID}); utils.HandleError(
		w,
		r,
		http.StatusInternalServerError,
		err,
		"error during encoding response",
	) {
		return
	} else {
		if _, err := w.Write(productIDJSON); utils.HandleError(
			w,
			r,
			http.StatusInternalServerError,
			err,
			"error writing response",
		) {
			return
		}
	}
}
