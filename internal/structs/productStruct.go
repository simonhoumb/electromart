package structs

import (
	"fmt"
	"go/types"
	"reflect"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	BrandID     string  `json:"brandID"`
	CategoryID  string  `json:"categoryID"`
	Description *string `json:"description"`
	QtyInStock  int     `json:"qtyInStock"`
	Price       float64 `json:"price"`
}

type ProductDiscount struct {
	DiscountID string `json:"discountID"`
	ProductID  string `json:"productID"`
}

type CreateProductResponse struct {
	ID string `json:"id"`
}

func (product Product) Validate() error {
	values := reflect.ValueOf(product)
	for i := 0; i < values.NumField(); i++ {
		switch values.Field(i).Interface().(type) {
		case int:
			if values.Field(i).Int() < 0 {
				return fmt.Errorf("field %s must be a positive integer", values.Type().Field(i).Name)
			}
		case types.Pointer:
			// Nil pointer are allowed
		default:
			if values.Field(i).IsZero() {
				return fmt.Errorf("field '%s' is invalid and/or required", values.Type().Field(i).Name)
			}
		}
	}
	return nil
}

func (product Product) ValidateNewProductRequest() error {
	values := reflect.ValueOf(product)
	for i := 0; i < values.NumField(); i++ {
		switch values.Field(i).Interface().(type) {
		case int:
			if values.Field(i).Int() < 0 {
				return fmt.Errorf("field %s must be a positive integer", values.Type().Field(i).Name)
			}
		case types.Pointer:
			// Nil pointer are allowed
		default:
			if values.Type().Field(i).Name == "ID" {
				// ID is not required for new products
			} else if values.Field(i).IsZero() {
				return fmt.Errorf("field '%s' is invalid and/or required", values.Type().Field(i).Name)
			}
		}
	}
	return nil
}
