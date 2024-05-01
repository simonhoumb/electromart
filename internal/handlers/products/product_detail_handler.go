package products

import (
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

}

func handleUpdateDetailRequest(w http.ResponseWriter, r *http.Request) {

}

func handleDeleteDetailRequest(w http.ResponseWriter, r *http.Request) {

}
