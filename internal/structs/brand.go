package structs

import "fmt"

type Brand struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateBrandResponse struct {
	ID string `json:"id"`
}

func (brand Brand) Validate() error {
	if brand.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}
