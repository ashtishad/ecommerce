package domain

import (
	"database/sql"
	"time"
)

type Category struct {
	CategoryID         int            `json:"categoryId"`
	CategoryUUID       string         `json:"categoryUuid"`
	ParentCategoryUUID sql.NullString `json:"parentCategoryUuid"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	Status             string         `json:"status"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	Level              int            `json:"level"`
	Subcategories      []Category     `json:"subcategories"`
}

func (c *Category) ToCategoryResponseDTO() *CategoryResponseDTO {
	subcategoriesDTO := make([]CategoryResponseDTO, len(c.Subcategories))
	for i, subcat := range c.Subcategories {
		subcategoriesDTO[i] = *subcat.ToCategoryResponseDTO()
	}

	parentUUID := ""
	if c.ParentCategoryUUID.Valid {
		parentUUID = c.ParentCategoryUUID.String
	}

	return &CategoryResponseDTO{
		CategoryUUID:       c.CategoryUUID,
		Name:               c.Name,
		Description:        c.Description,
		ParentCategoryUUID: parentUUID,
		Status:             c.Status,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
		Level:              c.Level,
		Subcategories:      subcategoriesDTO,
	}
}
