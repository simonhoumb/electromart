package utils

import (
	"Database_Project/structs"
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func RegisterUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		var creds structs.Credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 10)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		cartID := uuid.New().String()
		cartQuery := `INSERT INTO Cart (ID, TotalAmount) VALUES (?, ?)`
		_, err = db.Exec(cartQuery, cartID, 0)

		if err != nil {
			log.Fatalf("Failed to create cart: %v", err)
		}

		userID := uuid.New().String()
		query := `INSERT INTO User (ID, Username, Password, Email, FirstName, LastName, Phone, CartID) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(query, userID, creds.Username, string(hashedPassword), creds.Email, creds.FirstName, creds.LastName, creds.Phone, cartID)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			http.Error(w, "This username is taken", http.StatusConflict)
			return
		} else if err != nil {
			log.Printf("Failed to insert user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(structs.MessageResponse{Message: "User created successfully"})
	}
}
