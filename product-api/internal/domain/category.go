package domain

import "time"

type Category struct {
	CategoryID   int       `json:"categoryID"`
	CategoryUUID string    `json:"categoryUUID"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (c *Category) ToCategoryResponseDTO() *CategoryResponseDTO {
	return &CategoryResponseDTO{
		CategoryUUID:     c.CategoryUUID,
		Name:             c.Name,
		Description:      c.Description,
		Status:           c.Status,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		HasSubcategories: false,
		SubCategories:    nil,
	}
}
