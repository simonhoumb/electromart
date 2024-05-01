package structs

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     int64  `json:"phone"`
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
