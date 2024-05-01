package users

import (
	"Database_Project/internal/db"
	"Database_Project/internal/session"
	"Database_Project/internal/structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ProfileHandler handles user profile.
func ProfileHandler(userDB *db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleProfileGetRequest(w, r, userDB)
		case http.MethodPatch:
			handleProfilePatchRequest(w, r, userDB)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleProfileGetRequest now uses the session to get the username
func handleProfileGetRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	session, _ := session.Store.Get(r, "user-session")

	// Check if the user is logged in by looking for username in session.
	username, ok := session.Values["username"].(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println("Username from Session: ", username)
	user, err := userDB.GetUser(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond back with user details, make sure to not include sensitive details like password.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleProfilePatchRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	var user structs.ActiveUser

	userSession, _ := session.Store.Get(r, "user-session")
	if userSession.Values["username"] == nil {
		http.Error(w, "Unauthorized, please log in", http.StatusUnauthorized)
		return
	}

	// Decoding the user update request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Bad Request - Error Decoding JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if username := userSession.Values["username"]; username != user.Username {
		http.Error(w, "Permission Denied", http.StatusForbidden)
		return
	}

	err = userDB.UpdateUser(user)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	// Replying with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(structs.MessageResponse{Message: "User updated successfully"})
}
