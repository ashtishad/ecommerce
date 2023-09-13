package domain

import "time"

type CategoryResponseDTO struct {
	CategoryUUID       string                `json:"categoryUuid"`
	ParentCategoryUUID string                `json:"parentCategoryUuid"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	Status             string                `json:"status"`
	CreatedAt          time.Time             `json:"createdAt"`
	UpdatedAt          time.Time             `json:"updatedAt"`
	Level              int                   `json:"level"`
	Subcategories      []CategoryResponseDTO `json:"subcategories"`
}

type NewCategoryRequestDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
