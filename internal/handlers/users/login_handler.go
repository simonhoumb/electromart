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

// LoginHandler handles the login endpoint
func LoginHandler(userDB *db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleLoginPostRequest(w, r, userDB)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleLoginPostRequest handles the POST request to the login endpoint
func handleLoginPostRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	var loginRequest structs.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := userDB.GetUser(loginRequest.Username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	valid, err := userDB.CheckLogin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !valid {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create a new session and save the userID
	session, _ := session.Store.Get(r, "user-session") // get/create a session
	session.Values["userID"] = user.ID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"userID": user.ID})
}

// CheckLoginHandler checks whether user is logged in or not.
func CheckLoginHandler(userDB *db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := session.Store.Get(r, "user-session") // Call Store.Get() directly
		if err != nil {
			http.Error(w, "Unable to get session", http.StatusInternalServerError)
			log.Printf("Error fetching session: %v", err)
			return
		}
		userID := session.Values["userID"]
		if userID == nil {
			// Send response indicating user is not logged in
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"logged_in": false}`))
		} else {
			// Fetch username using userID
			user, err := userDB.GetUserByID(userID.(string))
			if err != nil {
				http.Error(w, "Error fetching user", http.StatusInternalServerError)
				log.Printf("Error fetching user: %v", err)
				return
			}

			// Send response indicating user is logged in along with their username
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"logged_in": true, "username": "` + user.Username + `"}`))
		}
	}
}

// LogoutHandler logs out the user.
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := session.Store.Get(r, "user-session")
		if err != nil {
			http.Error(w, "Unable to get session", http.StatusInternalServerError)
			log.Printf("Error fetching session: %v", err)
			return
		}

		// Clear session
		session.Values = make(map[interface{}]interface{})
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Unable to save session", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "Logout successful")
	}
}
