package structs

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	BrandID     string  `json:"brandID"`
	CategoryID  string  `json:"categoryID"`
	Description *string `json:"description"`
	QtyInStock  int     `json:"qtyInStock"`
	Price       float64 `json:"price"`
}
