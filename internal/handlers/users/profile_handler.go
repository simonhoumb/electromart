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

// handleProfilePatchRequest handles the PATCH request to update user profile.
// handleProfilePatchRequest handles the PATCH request to update user profile.
func handleProfilePatchRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	var user structs.ActiveUser

	// Retrieve user session
	userSession, err := session.Store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving user session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if the user is logged in by looking for username in session.
	username, ok := userSession.Values["username"].(string)
	if !ok {
		log.Println("No username found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode the user update request
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("User from request:", user.Username)
	log.Println("User from session:", username)

	// Ensure the user is updating their own profile
	if username != user.Username {
		log.Println("Permission Denied: Attempted to update another user's profile")
		http.Error(w, "Permission Denied", http.StatusForbidden)
		return
	}

	// Retrieve the existing user profile from the database
	existingUser, err := userDB.GetUser(user.Username)
	if err != nil {
		log.Printf("Error retrieving user profile: %v", err)
		http.Error(w, "Error retrieving user profile", http.StatusInternalServerError)
		return
	}

	// Update the user profile with the provided fields
	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}
	if user.Address.String != "" {
		existingUser.Address = user.Address
	}
	if user.PostCode.String != "" {
		existingUser.PostCode = user.PostCode
	}

	// Update the user profile in the database
	err = userDB.UpdateUserProfile(existingUser)
	if err != nil {
		log.Printf("Failed to update user profile: %v", err)
		http.Error(w, "Error updating user profile", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(structs.MessageResponse{Message: "User profile updated successfully"})
}
