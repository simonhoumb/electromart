package structs

import "database/sql"

type ActiveUser struct {
	ID        string         `json:"Id"`
	Username  string         `json:"Username"`
	Email     string         `json:"Email"`
	FirstName string         `json:"FirstName"`
	LastName  string         `json:"LastName"`
	Phone     string         `json:"Phone"`
	Address   sql.NullString `json:"Address,omitempty"`
	PostCode  sql.NullString `json:"PostCode,omitempty"`
	Password  string         `json:"Password"`
	CartID    string         `json:"CartID"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	CartID    string `json:"cartID"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
