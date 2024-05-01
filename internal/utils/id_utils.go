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
