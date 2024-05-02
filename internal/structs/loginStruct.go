package structs

type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     int    `json:"phone"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Username string `json:"username"`
}
