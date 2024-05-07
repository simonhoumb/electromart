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
