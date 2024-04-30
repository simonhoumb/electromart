package structs

type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     int    `json:"phone"`
}

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Username string `json:"username"`
}
