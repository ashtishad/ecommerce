package domain

import (
	"time"
)

type Product struct {
	ProductID      int
	ProductUUID    string
	CategoryID     int
	RootCategoryID int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PhoneAttribute struct {
	ProductID    int
	PriceInCents uint64
	SimType      SimType
	Storage      int
	Variant      VariantType
	Stock        int
}

type SoundEquipmentAttribute struct {
	ProductID    int
	PriceInCents uint64
	Codecs       []CodecType
	Color        string
	Stock        int
}

type WearablesAttribute struct {
	ProductID    int
	PriceInCents uint64
	StrapTypes   StrapType
	Color        string
	Stock        int
}

type SimType string

const (
	Single SimType = "Single"
	Dual   SimType = "Dual"
)

type CodecType string

const (
	SBC  CodecType = "SBC"
	AAC  CodecType = "AAC"
	Aptx CodecType = "Aptx"
	LDAC CodecType = "LDAC"
)

type StrapType string

const (
	Steel    StrapType = "Steel"
	Rubber   StrapType = "Rubber"
	Leather  StrapType = "Leather"
	Silicone StrapType = "Silicone"
)

type VariantType string

const (
	UAE      VariantType = "UAE"
	Official VariantType = "Official"
	USA      VariantType = "USA"
	China    VariantType = "China"
	UK       VariantType = "UK"
)
