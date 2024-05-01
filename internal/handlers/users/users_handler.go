package users

import (
	"Database_Project/internal/db"
	"net/http"
)

func UserHandler(userDB *db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleUserGetRequest(w, r)
			// If you have more methods to handle:
			// case http.MethodPost:
			//	handleUserPostRequest(w, r)
			// ... and so on for every HTTP method you need to process.
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleUserGetRequest(w http.ResponseWriter, r *http.Request) {
	// Your logic for handling the GET request on "/user" route goes here.
}
