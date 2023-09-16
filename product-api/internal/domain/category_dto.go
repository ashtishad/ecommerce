package domain

import "time"

type CategoryResponseDTO struct {
	CategoryUUID       string                `json:"categoryUuid"`
	ParentCategoryUUID string                `json:"parentCategoryUuid,omitempty"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	Status             string                `json:"status"`
	CreatedAt          time.Time             `json:"createdAt"`
	UpdatedAt          time.Time             `json:"updatedAt"`
	Level              int                   `json:"level,omitempty"`
	Subcategories      []CategoryResponseDTO `json:"subcategories,omitempty"`
}

type NewCategoryRequestDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
