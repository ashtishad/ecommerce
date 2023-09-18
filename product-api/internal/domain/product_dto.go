package domain

import (
	"time"
)

type Attributes any

type ProductDTO struct {
	ProductUUID string     `json:"productUuid"`
	Type        string     `json:"type"`
	SKU         string     `json:"sku"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Attributes  Attributes `json:"attributes"`
}

type PhoneAttributeDTO struct {
	PriceInCents uint64 `json:"priceInCents"`
	SimType      string `json:"simType"`
	Storage      string `json:"storage"`
	Variant      string `json:"variant"`
	Stock        int    `json:"stock"`
}

type SoundEquipmentAttributeDTO struct {
	PriceInCents uint64   `json:"priceInCents"`
	Codecs       []string `json:"codecs"`
	Color        string   `json:"color"`
	Stock        int      `json:"stock"`
}

type WearablesAttributeDTO struct {
	PriceInCents uint64   `json:"priceInCents"`
	StrapTypes   []string `json:"strapTypes"`
	Color        string   `json:"color"`
	Stock        int      `json:"stock"`
}
