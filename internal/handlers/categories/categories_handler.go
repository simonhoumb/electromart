package categories

import (
	"Database_Project/internal/db"
	"Database_Project/internal/structs"
	"Database_Project/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// Implemented methods for the endpoint
var categoriesImplementedMethods = []string{
	http.MethodGet,
	http.MethodPost,
}

/*
HandleCategories for the /categories endpoint.
*/
func HandleCategories(w http.ResponseWriter, r *http.Request) {
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
				categoriesImplementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleGetAllRequest(w http.ResponseWriter, r *http.Request) {
	// Get all categories
	categories, err := db.GetAllCategories()
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error getting categories from database") {
		return
	}

	// Return the categories
	if categoriesJSON, err := json.Marshal(categories); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
		return
	} else {
		if _, err := w.Write(categoriesJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
			return
		}
	}
}

func handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	var category structs.Category

	if err := json.NewDecoder(r.Body).Decode(&category); utils.HandleError(w, r, http.StatusBadRequest, err, "error during decoding request") {
		return
	}

	if err := category.Validate(); utils.HandleError(w, r, http.StatusBadRequest, err, "invalid request json, check documentation") {
		return
	}

	// Create the category
	err := db.AddCategory(category)
	if utils.HandleError(w, r, http.StatusInternalServerError, err, "error adding category to database") {
		return
	}

	// No content to return
	w.WriteHeader(http.StatusNoContent)
}
