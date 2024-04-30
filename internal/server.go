package server

import (
<<<<<<< HEAD
	"Database_Project/db"
	"Database_Project/utils"
=======
	"Database_Project/internal/db"
	utils2 "Database_Project/internal/utils"
>>>>>>> b60bf6058593cea444e1a5258bb777f2445535de
	"log"
	"net/http"
)

// Start
/*
Start the server on the port specified in the environment variable PORT. If PORT is not set, the default port 8080 is used.
*/
func Start() {

	// Get the port from the environment variable, or use the default port
<<<<<<< HEAD
	port := utils.GetPort()
=======
	port := utils2.GetPort()
>>>>>>> b60bf6058593cea444e1a5258bb777f2445535de

	// Using mux to handle /'s and parameters
	mux := http.NewServeMux()

	db.Client = db.Connect()
	defer db.Client.Close()

	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/index.html")
		},
	)
<<<<<<< HEAD
	mux.HandleFunc("/login", utils.CheckLogin(db.Client))
	mux.HandleFunc("/logout", utils.LogoutUser(db.Client))
	mux.HandleFunc(
		"/cart", func(w http.ResponseWriter, r *http.Request) {
			utils.GetCartItems(w, r, db.Client)
=======
	mux.HandleFunc("/login", utils2.CheckLogin(db.Client))
	mux.HandleFunc("/logout", utils2.LogoutUser(db.Client))
	mux.HandleFunc(
		"/cart", func(w http.ResponseWriter, r *http.Request) {
			utils2.GetCartItems(w, r, db.Client)
>>>>>>> b60bf6058593cea444e1a5258bb777f2445535de
		},
	)
	mux.HandleFunc(
		"/profile", func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
			utils.GetUserProfile(w, r, db.Client)
		},
	)
	mux.HandleFunc("/api/categories", db.GetCategoriesHandler(db.Client))
	mux.HandleFunc("/register", utils.RegisterUser(db.Client))
=======
			utils2.GetUserProfile(w, r, db.Client)
		},
	)
	mux.HandleFunc("/api/categories", db.GetCategoriesHandler(db.Client))
	mux.HandleFunc("/register", utils2.RegisterUser(db.Client))
>>>>>>> b60bf6058593cea444e1a5258bb777f2445535de
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc(
		"/loginPage", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/Login.html")
		},
	)

	mux.HandleFunc(
		"/registerPage", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/register.html")
		},
	)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
