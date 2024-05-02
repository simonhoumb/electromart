package structs

import "time"

type Order struct {
	ID                string    `json:"id"`
	UserID            string    `json:"userID"`
	TotalAmount       float64   `json:"totalAmount"`
	OrderDate         string    `json:"orderDate"`
	ShippedDate       time.Time `json:"shippedDate"`
	EstimatedDelivery time.Time `json:"estimatedDelivery"`
	Status            string    `json:"status"`
	Comments          *string   `json:"comments"`
}
