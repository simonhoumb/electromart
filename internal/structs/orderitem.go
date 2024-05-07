package structs

type OrderItem struct {
	ProductID      string  `json:"productID"`
	ProductOrderID string  `json:"productOrderID"`
	Quantity       int     `json:"quantity"`
	SubTotal       float64 `json:"subTotal"`
}
