package structs

import "fmt"

type Category struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateCategoryResponse struct {
	ID string `json:"id"`
}

func (category Category) Validate() error {
	if category.ID != "" {
		return fmt.Errorf("field 'id' is invalid")
	}
	if category.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}

func (category Category) ValidateNewCategoryRequest() error {
	if category.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	return nil
}
