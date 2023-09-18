package domain

import (
	"time"
)

type Attributes any

type ProductDTO struct {
	ProductUUID string     `json:"productUuid"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	SKU         string     `json:"sku"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Attributes  Attributes `json:"attributes"`

	// consider adding category_uuid and root_category_uuid or name later
}

type PhoneAttributeDTO struct {
	Price   string `json:"price"` // display price string, priceInCents/100
	SimType string `json:"simType"`
	Storage int    `json:"storage"`
	Variant string `json:"variant"`
	Stock   int    `json:"stock"`
}

type SoundEquipmentAttributeDTO struct {
	Price  string   `json:"price"` // display price string, priceInCents/100
	Codecs []string `json:"codecs"`
	Color  string   `json:"color"`
	Stock  int      `json:"stock"`
}

type WearablesAttributeDTO struct {
	Price      string   `json:"price"` // display price string, priceInCents/100
	StrapTypes []string `json:"strapTypes"`
	Color      string   `json:"color"`
	Stock      int      `json:"stock"`
}
