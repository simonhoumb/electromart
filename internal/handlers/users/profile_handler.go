package users

import (
	"Database_Project/internal/db"
	"Database_Project/internal/session"
	"Database_Project/internal/structs"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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
	// 1. Get the session
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// 2. Retrieve the userID from the session
	userIDValue := session.Values["userID"].(string)

	// 4. Fetch user by ID
	user, err := userDB.GetUserByID(userIDValue)
	if err != nil {
		http.Error(w, "Error retrieving user profile", http.StatusInternalServerError)
		return
	}

	// 5. Mask the password
	user.Password = ""

	respondWithJSON(w, http.StatusOK, user)
}

// handleProfilePatchRequest handles the PATCH request to update user profile.
func handleProfilePatchRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	// 1. Get the session
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// 2. Retrieve the userID from the session
	userIDValue := session.Values["userID"].(string)

	// 3. Decode the request body
	var user structs.ActiveUser
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 4. Update the user profile
	existingUser, err := userDB.GetUserByID(userIDValue)
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

// handleProfileDeleteRequest handles deleting a user after password confirmation.
// handleProfileDeleteRequest handles deleting a user after password confirmation.
func handleProfileDeleteRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	// 1. Get the session
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Retrieve the userID from the session
	userIDValue := session.Values["userID"].(string)

	// 3. Decode the request body to get the password confirmation
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 4. Check if password confirmation is provided
	passwordConfirmation, ok := requestBody["passwordConfirmation"]
	if !ok || passwordConfirmation == "" {
		http.Error(w, "Password confirmation is required", http.StatusBadRequest)
		return
	}

	// 5. Retrieve user information from the database using userID
	user, err := userDB.GetUserByID(userIDValue)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 6. Verify password confirmation
	if !verifyPasswordHash(passwordConfirmation, user.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// 7. Password confirmed, proceed with deleting the user
	err = deleteUser(userIDValue, userDB) // Use userID.String() for the query
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	// 8. Clear the session (log the user out)
	clearSession(w, r)

	// 9. Respond with success
	w.WriteHeader(http.StatusNoContent)
}

// verifyPasswordHash verifies if the provided password matches the hashed password stored in the database.
func verifyPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func getUsernameFromSession(r *http.Request) (string, error) {
	session, err := session.Store.Get(r, "user-session")
	if err != nil {
		return "", err
	}
	username, ok := session.Values["userID"].(string)
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

func deleteUser(userID string, userDB *db.UserDB) error {
	err := userDB.DeleteUser(userID) // Assuming DeleteUser now takes a uuid.UUID
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
