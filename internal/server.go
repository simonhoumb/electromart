package server

import (
	"Database_Project/internal/constants"
	"Database_Project/internal/db"
	"Database_Project/internal/handlers/brands"
	"Database_Project/internal/handlers/categories"
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

	db.Client = db.OpenDatabaseConnection()
	defer db.Client.Close()

	// API endpoints
	// Handle the products endpoint
	mux.HandleFunc(constants.ProductsPath, products.HandleProducts)
	mux.HandleFunc(constants.ProductsPath+"{id}", products.HandleProductDetail)
	mux.HandleFunc(constants.ProductsPath+"search/{query}", products.HandleQueryProducts)

	// Handle the categories endpoint
	mux.HandleFunc(constants.CategoriesPath, categories.HandleCategories)
	mux.HandleFunc(constants.CategoriesPath+"{name}", categories.HandleCategoryDetail)

	// Handle the brands endpoint
	mux.HandleFunc(constants.BrandsPath, brands.HandleBrands)
	mux.HandleFunc(constants.BrandsPath+"{name}", brands.HandleBrandDetail)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	mux.HandleFunc("/api/check_login", users.CheckLoginHandler(userDB))
	mux.HandleFunc("/api/logout", users.LogoutHandler())
	mux.HandleFunc("/api/login", users.LoginHandler(userDB))
	mux.HandleFunc("/api/register", users.RegistrationHandler(userDB))
	mux.HandleFunc("/api/profile", session.CheckSession(users.ProfileHandler(userDB)))
	mux.HandleFunc("/api/change_password", session.CheckSession(users.ChangePasswordHandler(userDB)))

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/Login.html")
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/register.html")
	})

	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/profile.html")
	})

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/product.html")
	})

	port := utils.GetPort()

	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getUserSession(r *http.Request) (*sessions.Session, error) {
	return session.Store.Get(r, "user-session") // get/create a session
}
