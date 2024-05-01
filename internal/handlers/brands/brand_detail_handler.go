package brands

import (
	"fmt"
	"net/http"
)

var detailImplementedMethods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodDelete,
}

func HandleBrandDetail(w http.ResponseWriter, r *http.Request) {
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
	//id, err := utils.GetIDFromRequest(r)
	//if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
	//	return
	//}
	//
	//// Get the brand with the given ID
	//brand, err := db.GetBrandByID(db.Client, id)
	//if utils.HandleError(w, r, http.StatusInternalServerError, err, "Error getting brands from database") {
	//	return
	//}
	//
	//// Return the brand
	//if marshalledBrand, err := json.Marshal(brand); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error during encoding response") {
	//	return
	//} else {
	//	if _, err := w.Write(marshalledBrand); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error writing response") {
	//		return
	//	}
	//}
}

func handleUpdateDetailRequest(w http.ResponseWriter, r *http.Request) {
	//id, err := utils.GetIDFromRequest(r)
	//if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
	//	return
	//}
	//
	//// Decode the request body into a brand
	//var updatedBrand structs.Brand
	//if err := json.NewDecoder(r.Body).Decode(&updatedBrand); utils.HandleError(w, r, http.StatusBadRequest, err, "Error decoding request body") {
	//	return
	//}
	//
	//if updatedBrand.ID != id {
	//	utils.HandleError(w, r, http.StatusBadRequest, fmt.Errorf("ID in request body does not match ID in URL"), "ID in request body does not match ID in URL")
	//	return
	//}
	//
	//// Update the brand with the given ID
	//if err := db.UpdateBrand(db.Client, updatedBrand); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error updating brand in database") {
	//	return
	//}
	//
	//// Return no content
	//w.WriteHeader(http.StatusNoContent)
}

func handleDeleteDetailRequest(w http.ResponseWriter, r *http.Request) {
	//id, err := utils.GetIDFromRequest(r)
	//if utils.HandleError(w, r, http.StatusBadRequest, err, "Error getting ID from request") {
	//	return
	//}
	//
	//// Get the brand with the given ID
	//if err := db.DeleteBrandByID(db.Client, id); utils.HandleError(w, r, http.StatusInternalServerError, err, "Error deleting brand from database") {
	//	return
	//}
	//
	//// Return no content
	//w.WriteHeader(http.StatusNoContent)
}
