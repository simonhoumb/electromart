package server

import (
	"Database_Project/internal/constants"
	"Database_Project/internal/db"
	"Database_Project/internal/handlers/categories"
	"Database_Project/internal/handlers/products"
	"Database_Project/internal/utils"
	"log"
	"net/http"
)

// Start
/*
Start the server on the port specified in the environment variable PORT. If PORT is not set, the default port 8080 is used.
*/
func Start() {

	// Get the port from the environment variable, or use the default port
	port := utils.GetPort()

	// Using mux to handle /'s and parameters
	mux := http.NewServeMux()

	db.Client = db.OpenDatabaseConnection()

	defer db.Client.Close()

	// Handle the products endpoint
	mux.HandleFunc(constants.ProductsPath, products.HandleProducts)
	mux.HandleFunc(constants.ProductsPath+"{id}", products.HandleProductDetail)

	// Handle the categories endpoint
	mux.HandleFunc(constants.CategoriesPath, categories.Handler)

	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/index.html")
		},
	)

	mux.HandleFunc("/login", utils.CheckLogin(db.Client))
	mux.HandleFunc("/logout", utils.LogoutUser(db.Client))
	mux.HandleFunc(
		"/cart", func(w http.ResponseWriter, r *http.Request) {
			utils.GetCartItems(w, r, db.Client)

		},
	)

	mux.HandleFunc(
		"/profile", func(w http.ResponseWriter, r *http.Request) {
			utils.GetUserProfile(w, r, db.Client)
		},
	)

	mux.HandleFunc("/register", utils.RegisterUser(db.Client))

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
