package structs

import "fmt"

type BrandName struct {
	Name string `json:"name"`
}

type Brand struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (brand Brand) Validate() error {
	if brand.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}
