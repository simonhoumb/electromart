package utils

import (
	"Database_Project/internal/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-unique-secret"))

func CheckLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		var creds structs.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			fmt.Println("Error in JSON decoding:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		fmt.Println("Received creds: ", creds)

		row := db.QueryRow("SELECT password FROM User WHERE username = ?", creds.Username)

		var hashedPassword string
		err = row.Scan(&hashedPassword)

		if err == sql.ErrNoRows {
			fmt.Println("Username not found")
			http.Error(w, "Username not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			fmt.Println("internal Error when fetching row:", err)
			http.Error(w, "internal Error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
			fmt.Println("Password check failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println("Login check succeeded")

		username := creds.Username
		type Response struct {
			Username string `json:"username"`
		}

		resp := Response{Username: username}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Error in JSON encoding:", err)
			http.Error(w, "internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func LogoutUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		session.Options = &sessions.Options{
			MaxAge: -1,
		}
		session.Save(r, w)
		http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
	}
}

func GetUserProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Write([]byte("User Profile: Not Implemented Yet"))
}

func GetCartItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Write([]byte("Cart Items: Not Implemented Yet"))
}
