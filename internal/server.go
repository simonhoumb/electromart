package server

import (
	"Database_Project/internal/constants"
	"Database_Project/internal/db"
	"Database_Project/internal/handlers/products"
	"Database_Project/internal/handlers/users"
	"Database_Project/internal/session"
	"Database_Project/internal/utils"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func Start() {
	// Using mux to handle /'s and parameters
	mux := http.NewServeMux()

	userDB := &db.UserDB{Client: db.OpenDatabaseConnection()}
	defer userDB.Client.Close()

	mux.HandleFunc(constants.ProductsPath, products.Handler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	mux.HandleFunc("/api/check_login", users.CheckLoginHandler(userDB))
	mux.HandleFunc("/api/logout", users.LogoutHandler())
	mux.HandleFunc("/api/login", users.LoginHandler(userDB))

	mux.HandleFunc("/register", utils.RegisterUser(userDB.Client))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("/loginPage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/Login.html")
	})

	mux.HandleFunc("/registerPage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/register.html")
	})

	port := utils.GetPort()

	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getUserSession(r *http.Request) (*sessions.Session, error) {
	return session.Store.Get(r, "user-session") // get/create a session
}
