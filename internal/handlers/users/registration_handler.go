package users

import (
	"Database_Project/internal/db"
	"Database_Project/internal/structs"
	"encoding/json"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func RegistrationHandler(userDB *db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleRegistrationPostRequest(w, r, userDB)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// RegistrationHandler handles user registration.
func handleRegistrationPostRequest(w http.ResponseWriter, r *http.Request, userDB *db.UserDB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect method", http.StatusNotFound)
		return
	}

	var creds structs.ActiveUser

	// Decoding the user registration request
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("Bad Request - Error Decoding JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 10)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error while encrypting password")
		return
	}

	// Inserting data into Cart
	cartID, err := userDB.CreateUserCart()
	if err != nil {
		log.Printf("Failed to create cart: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user cart")
		return
	}

	// Inserting data into User
	userID := uuid.New().String()
	err = userDB.RegisterUser(userID, creds.Username, string(hashedPassword), creds.Email, creds.FirstName, creds.LastName, creds.Phone, cartID, creds.Address, creds.PostCode)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == 1062 {
			respondWithError(w, http.StatusConflict, "This username is taken")
			return
		}
		log.Printf("Failed to register user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error registering user")
		return
	}

	// Replying with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(structs.MessageResponse{Message: "User created successfully"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	response := structs.MessageResponse{Message: message}
	json.NewEncoder(w).Encode(response)
}
