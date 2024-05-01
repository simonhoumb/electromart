package products

import (
	"Database_Project/internal/db"
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

var detailImplementedMethods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

func HandleProductDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// Switch on the HTTP request method
	switch r.Method {
	case http.MethodGet:
		handleGetDetailRequest(w, r)

	case http.MethodPut:
		handleUpdateDetailRequest(w, r)

	case http.MethodDelete:
		handleDeleteDetailRequest(w, r)

	default:
		// If the method is not implemented, return an error with the allowed methods
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				detailImplementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleGetDetailRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromRequest(r)
	if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
		return
	}

	// Get the product with the given ID
	products, err := db.GetAllProducts(db.Client)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "Error getting products from database") {
		return
	}

	// Take first product from the list TODO: Change this to get the product with the given ID
	product := products[0]
	product.ID = product.ID + id

	// Return the product
	if marshalledProduct, err := json.Marshal(product); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error during encoding response") {
		return
	} else {
		if _, err := w.Write(marshalledProduct); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error writing response") {
			return
		}
	}
}

func handleUpdateDetailRequest(w http.ResponseWriter, r *http.Request) {

}

func handleDeleteDetailRequest(w http.ResponseWriter, r *http.Request) {

}
