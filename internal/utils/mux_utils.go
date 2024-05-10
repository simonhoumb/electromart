package utils

import (
	"errors"
	"net/http"
)

// GetIDFromRequest Get the ID from the mux request
func GetIDFromRequest(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if id == "" {
		return "", errors.New("ID not found in request")
	}

	return id, nil
}

// GetNameFromRequest Get the name from the mux request
func GetNameFromRequest(r *http.Request) (string, error) {
	name := r.PathValue("name")
	if name == "" {
		return "", errors.New("name not found in request")
	}

	return name, nil
}

// GetQueryFromRequest Get the query from the mux request
func GetQueryFromRequest(r *http.Request) (string, error) {
	query := r.PathValue("query")
	if query == "" {
		return "", errors.New("query not found in request")
	}

	return query, nil
}

// GetCategoryFromRequest Get the category from the mux request
func GetCategoryFromRequest(r *http.Request) (string, error) {
	category := r.PathValue("category")
	if category == "" {
		return "", errors.New("category not found in request")
	}

	return category, nil
}

// GetBrandFromRequest Get the brand from the mux request
func GetBrandFromRequest(r *http.Request) (string, error) {
	brand := r.PathValue("brand")
	if brand == "" {
		return "", errors.New("brand not found in request")
	}

	return brand, nil
}
