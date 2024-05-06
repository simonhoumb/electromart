package structs

type OrderItem struct {
	ID        string `json:"id"`
	OrderID   string `json:"orderID"`
	ProductID string `json:"productID"`
	Quantity  int    `json:"quantity"`
}
