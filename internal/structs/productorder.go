package structs

import "time"

type ProductOrder struct {
	ID                string     `json:"id"`
	UserAccountID     string     `json:"userID"`
	OrderDate         time.Time  `json:"orderDate"`
	ShippedDate       *time.Time `json:"shippedDate"`
	EstimatedDelivery *time.Time `json:"estimatedDelivery"`
	DeliveryFee       float64    `json:"deliveryFee"`
	Status            string     `json:"status"`
	Comments          *string    `json:"comments"`
}
