package domain

import "time"

type CategoryResponseDTO struct {
	CategoryUUID     string             `json:"categoryUuid"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Status           string             `json:"status"`
	CreatedAt        time.Time          `json:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	HasSubcategories bool               `json:"hasSubcategories"`
	SubCategories    []SubCategoryBrief `json:"subCategories"`
}

type SubCategoryBrief struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type NewCategoryRequestDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
