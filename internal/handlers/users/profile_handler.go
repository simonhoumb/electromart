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
		case http.MethodDelete:
			handleProfileDeleteRequest(w, r, userDB)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleProfileGetRequest now uses the session to get the username
func handleProfileGetRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := userDB.GetUser(username)
	if err != nil {
		http.Error(w, "Error retrieving user profile", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// handleProfilePatchRequest handles the PATCH request to update user profile.
func handleProfilePatchRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user structs.ActiveUser
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if username != user.Username {
		http.Error(w, "Permission Denied", http.StatusForbidden)
		return
	}

	existingUser, err := userDB.GetUser(username)
	if err != nil {
		http.Error(w, "Error retrieving user profile", http.StatusInternalServerError)
		return
	}

	err = updateUserProfile(existingUser, user, userDB)
	if err != nil {
		http.Error(w, "Error updating user profile", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, structs.MessageResponse{Message: "User profile updated successfully"})
}

// handleProfileDeleteRequest handles deleting a user.
func handleProfileDeleteRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = deleteUser(username, userDB)
	if err != nil {
		log.Printf("Error deleting user: %v", err) // Log the specific error
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	clearSession(w, r)
	respondWithText(w, http.StatusNoContent, "User deleted successfully")
}

func getUsernameFromSession(r *http.Request) (string, error) {
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		return "", err
	}
	username, ok := session.Values["username"].(string)
	if !ok {
		return "", fmt.Errorf("no username found in session")
	}
	return username, nil
}

func updateUserProfile(existingUser, user structs.ActiveUser, userDB *db.UserDB) error {
	// Update the user profile with the provided fields
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}
	if user.Phone != "" {
		existingUser.Phone = user.Phone
	}
	if user.Address.String != "" {
		existingUser.Address = user.Address
	}
	if user.PostCode.String != "" {
		existingUser.PostCode = user.PostCode
	}

	err := userDB.UpdateUserProfile(existingUser)
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(username string, userDB *db.UserDB) error {
	err := userDB.DeleteUser(username)
	if err != nil {
		return err
	}
	return nil
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving user session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error clearing user session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondWithText(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}
