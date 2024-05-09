package users

import (
	"Database_Project/internal/db"
	"Database_Project/internal/structs"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// ChangePasswordHandler handles the logic for changing a user's password.
func ChangePasswordHandler(userDB *db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handleChangePassword(w, r, userDB)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleChangePassword(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	// Retrieve username from session
	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var changePasswordRequest structs.ChangePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&changePasswordRequest)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Verify old password
	isValid, err := userDB.CheckLogin(username, changePasswordRequest.OldPassword)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !isValid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update password in the database
	err = userDB.UpdatePassword(username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	respondWithText(w, http.StatusOK, "Password changed successfully")
}
