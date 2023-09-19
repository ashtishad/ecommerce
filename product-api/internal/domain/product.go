package domain

import "time"

type Product struct {
	ProductID    int
	ProductUUID  string
	Name         string
	Color        string
	CategoryUUID string
	ProductType  ProductType // example: "Phone", "SoundEquipment" -> root category level 0
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Attributes   []ProductAttribute
}

type ProductAttribute interface {
	SetAttributes() ProductAttribute
}

type PhoneAttribute struct {
	ProductID    int
	PriceInCents uint64
	SimType      SimType
	Storage      StorageType
	Variant      VariantType
	Stock        int
}

func (pa *PhoneAttribute) SetAttributes() ProductAttribute {
	return &PhoneAttribute{
		ProductID:    pa.ProductID,
		PriceInCents: pa.PriceInCents,
		SimType:      pa.SimType,
		Storage:      pa.Storage,
		Variant:      pa.Variant,
		Stock:        pa.Stock,
	}
}

type SoundEquipmentAttribute struct {
	ProductID    int
	PriceInCents uint64
	Codecs       []CodecType
	Stock        int
}

func (sa *SoundEquipmentAttribute) SetAttributes() ProductAttribute {
	return &SoundEquipmentAttribute{
		ProductID:    sa.ProductID,
		PriceInCents: sa.PriceInCents,
		Codecs:       sa.Codecs,
		Stock:        sa.Stock,
	}
}

type WearablesAttribute struct {
	ProductID    int
	PriceInCents uint64
	StrapTypes   StrapType
	Stock        int
}

func (wa *WearablesAttribute) SetAttributes() ProductAttribute {
	return &WearablesAttribute{
		ProductID:    wa.ProductID,
		PriceInCents: wa.PriceInCents,
		StrapTypes:   wa.StrapTypes,
		Stock:        wa.Stock,
	}
}

// Pre-defined types

type ProductType string

const (
	ProductTypePhone          ProductType = "Phone"
	ProductTypeSoundEquipment ProductType = "SoundEquipment"
	ProductTypeWearable       ProductType = "Wearable"
)

type SimType string

const (
	SimTypeSingle SimType = "Single"
	SimTypeDual   SimType = "Dual"
	SimTypeEsim   SimType = "eSim"
)

type StorageType int

const (
	Storage128GB StorageType = 128
	Storage256GB StorageType = 256
	Storage512GB StorageType = 512
	Storage1TB   StorageType = 1024
)

type CodecType string

const (
	CodecTypeSBC  CodecType = "SBC"
	CodecTypeAAC  CodecType = "AAC"
	CodecTypeAptx CodecType = "Aptx"
	CodecTypeLDAC CodecType = "LDAC"
)

type StrapType string

const (
	StrapTypeSteel    StrapType = "Steel"
	StrapTypeRubber   StrapType = "Rubber"
	StrapTypeLeather  StrapType = "Leather"
	StrapTypeSilicone StrapType = "Silicone"
)

type VariantType string

const (
	VariantTypeUAE      VariantType = "UAE"
	VariantTypeOfficial VariantType = "Official"
	VariantTypeUSA      VariantType = "USA"
	VariantTypeChina    VariantType = "China"
)
