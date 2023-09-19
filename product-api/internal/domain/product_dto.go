package domain

import (
	"time"

	"github.com/ashtishad/ecommerce/lib"
)

type ProductResponseDTO struct {
	// common fields for a product
	ProductUUID  string      `json:"productUuid"`
	ProductName  string      `json:"productName"`
	ProductType  ProductType `json:"productType"`
	CreatedAt    time.Time   `json:"createdAt"`
	CategoryUUID string      `json:"categoryUuid"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	DisplayPrice string      `json:"displayPrice"`
	Stock        int         `json:"stock"`
	Color        string      `json:"color"`

	// Phone
	SimType SimType     `json:"simType,omitempty"`
	Storage StorageType `json:"storage,omitempty"`
	Variant VariantType `json:"variant,omitempty"`

	// SoundEquipment
	Codecs []CodecType `json:"codecs,omitempty"`

	// Wearable
	StrapTypes StrapType `json:"strapTypes,omitempty"`
}

func (p Product) ToProductResponseDTO() *ProductResponseDTO {
	res := ProductResponseDTO{
		ProductUUID:  p.ProductUUID,
		ProductName:  p.Name,
		Color:        p.Color,
		CategoryUUID: p.CategoryUUID,
		ProductType:  p.ProductType,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	switch p.ProductType {
	case ProductTypePhone:
		var pa PhoneAttribute
		pa.SetAttributes()
		res.SimType = pa.SimType
		res.Storage = pa.Storage
		res.Variant = pa.Variant
		res.Stock = pa.Stock
		res.DisplayPrice = lib.DisplayPrice(int(pa.PriceInCents))

	case ProductTypeSoundEquipment:
		var sa SoundEquipmentAttribute
		sa.SetAttributes()
		res.Codecs = sa.Codecs
		res.Stock = sa.Stock
		res.DisplayPrice = lib.DisplayPrice(int(sa.PriceInCents))

	case ProductTypeWearable:
		var wa WearablesAttribute
		wa.SetAttributes()
		res.StrapTypes = wa.StrapTypes
		res.Stock = wa.Stock
		res.DisplayPrice = lib.DisplayPrice(int(wa.PriceInCents))
	}

	return &res
}
