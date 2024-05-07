package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

// Store is exported and can be accessed from other packages
var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte("your-very-secret-key"))
	Store.Options = &sessions.Options{
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		Path:     "/", // Add this line
		SameSite: http.SameSiteStrictMode,
	}
}

// In your session package
func CheckSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "user-session")
		if err != nil {
			log.Println("Error getting session:", err)
			http.Error(w, "Unable to get session", http.StatusInternalServerError)
			return
		}
		if session.Values["username"] == nil {
			log.Println("Unauthorized access attempt")
			fmt.Fprintln(w, "false")
			return
		}
		next.ServeHTTP(w, r)
	}
}
