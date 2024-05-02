package structs

import "fmt"

type Brand struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateBrandResponse struct {
	ID string `json:"id"`
}

func (brand Brand) Validate() error {
	if brand.ID != "" {
		return fmt.Errorf("field 'id' is invalid")
	}
	if brand.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}

func (brand Brand) ValidateNewBrandRequest() error {
	if brand.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}
