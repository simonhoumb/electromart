package categories

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	// Get all categories
	// categories := db.somethingsomething()
	categories := []string{"category1", "category2", "category3"}

	// Marshal the products into a JSON object
	productsJSON, err := json.Marshal(categories)
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
