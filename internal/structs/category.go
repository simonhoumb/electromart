package structs

import "fmt"

type Category struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (category Category) Validate() error {
	if category.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}
