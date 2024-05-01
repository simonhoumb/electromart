package utils

import (
	"log"
	"net/http"
)

/*
HandleError handles an error by sending an HTTP error response with the given status code and error message.

Example usage in an HTTP handler:

For a function that returns an error:

	if err := funcThatReturnsError; utils.HandleError(w, r, http.StatusBadRequest, err, "error message to user") {
		return
	}

For a function that returns an error along with another value (e.g. a response):

	response, err := funcThatReturnsErrorAndResponse()
	if utils.HandleError(w, r, http.StatusBadRequest, err, "error message to user") {
		return
	}

See product_handler.go for an example of how to use this function. TODO: Remove this line
*/
func HandleError(w http.ResponseWriter, r *http.Request, status int, err error, errorMessage string) bool {
	if err != nil {
		http.Error(w, errorMessage, status)
		log.Printf("Error handled by HandleError:\n\tError: %v\n\tRequest: %v", err, r)
		return true
	}
	return false
}
