package brands // Implemented methods for the endpoint
import (
	"fmt"
	"net/http"
)

var brandsImplementedMethods = []string{
	http.MethodGet,
	http.MethodPost,
}

/*
HandleBrands for the /brands endpoint.
*/
func HandleBrands(w http.ResponseWriter, r *http.Request) {
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
				brandsImplementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

func handleGetAllRequest(w http.ResponseWriter, r *http.Request) {
	//// Get all brands
	//brands, err := db.GetAllBrands(db.Client)
	//if utils.HandleError(w, r, http.StatusInternalServerError, err, "error getting brands from database") {
	//	return
	//}
	//
	//// Return the brands
	//if brandsJSON, err := json.Marshal(brands); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
	//	return
	//} else {
	//	if _, err := w.Write(brandsJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
	//		return
	//	}
	//}
}

func handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	//var brand structs.Brand
	//
	//if err := json.NewDecoder(r.Body).Decode(&brand); utils.HandleError(w, r, http.StatusBadRequest, err, "error during decoding request") {
	//	return
	//}
	//
	//if err := brand.ValidateNewBrandRequest(); utils.HandleError(w, r, http.StatusBadRequest, err, "invalid request json, check documentation") {
	//	return
	//}
	//
	//// Create the brand
	//brandID, err := db.AddBrand(db.Client, brand)
	//if utils.HandleError(w, r, http.StatusInternalServerError, err, "error adding brand to database") {
	//	return
	//}
	//
	//// Two above in one if statement
	//if brandIDJSON, err := json.Marshal(structs.CreateBrandResponse{ID: brandID}); utils.HandleError(w, r, http.StatusInternalServerError, err, "error during encoding response") {
	//	return
	//} else {
	//	if _, err := w.Write(brandIDJSON); utils.HandleError(w, r, http.StatusInternalServerError, err, "error writing response") {
	//		return
	//	}
	//}
}
