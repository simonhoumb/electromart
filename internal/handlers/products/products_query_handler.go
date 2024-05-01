package products

import (
	"Database_Project/internal/db"
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

var queryImplementedMethods = []string{
	http.MethodGet,
}

func HandleQueryProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleProductsQueryRequest(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				queryImplementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleProductsQueryRequest(w http.ResponseWriter, r *http.Request) {
	query, err := utils.GetQueryFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting query from request") {
		return
	}

	// Get the product with the given ID
	products, err := db.SearchProducts(db.Client, query)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "Error getting products from database") {
		return
	}

	// Return the product
	if marshalledProducts, err := json.Marshal(products); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error during encoding response") {
		return
	} else {
		if _, err := w.Write(marshalledProducts); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error writing response") {
			return
		}
	}
}
