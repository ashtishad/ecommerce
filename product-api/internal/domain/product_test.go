package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetPhoneObject(productID int, priceInCents uint64, simType SimType, storage StorageType, variant VariantType, stock int) *PhoneAttribute {
	return &PhoneAttribute{
		ProductID:    productID,
		PriceInCents: priceInCents,
		SimType:      simType,
		Storage:      storage,
		Variant:      variant,
		Stock:        stock,
	}
}

func TestPhoneAttributeSetAttributes(t *testing.T) {
	tests := []struct {
		name         string
		phoneAttr    *PhoneAttribute
		expectedAttr *PhoneAttribute
	}{
		{
			name:         "Test Case 1",
			phoneAttr:    GetPhoneObject(1, 10099, SimTypeDual, Storage128GB, VariantTypeOfficial, 10),
			expectedAttr: GetPhoneObject(1, 10099, SimTypeDual, Storage128GB, VariantTypeOfficial, 10),
		},
		{
			name:         "Test Case 2",
			phoneAttr:    GetPhoneObject(2, 20099, SimTypeSingle, Storage256GB, VariantTypeUAE, 20),
			expectedAttr: GetPhoneObject(2, 20099, SimTypeSingle, Storage256GB, VariantTypeUAE, 20),
		},
		{
			name:         "Test Case 3",
			phoneAttr:    GetPhoneObject(3, 145678, SimTypeSingle, Storage512GB, VariantTypeUSA, 20),
			expectedAttr: GetPhoneObject(3, 145678, SimTypeSingle, Storage512GB, VariantTypeUSA, 20),
		},
		{
			name:         "Test Case 4",
			phoneAttr:    GetPhoneObject(4, 12345, SimTypeDual, Storage1TB, VariantTypeChina, 10),
			expectedAttr: GetPhoneObject(4, 12345, SimTypeDual, Storage1TB, VariantTypeChina, 10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.phoneAttr.SetAttributes().(*PhoneAttribute)
			assert.Equal(t, *tt.expectedAttr, *result, fmt.Sprintf("%s: Expected and actual PhoneAttributes are not equal", tt.name))
		})
	}
}

func GetSoundEquipmentObject(productID int, priceInCents uint64, codecs []CodecType, stock int) *SoundEquipmentAttribute {
	return &SoundEquipmentAttribute{
		ProductID:    productID,
		PriceInCents: priceInCents,
		Codecs:       codecs,
		Stock:        stock,
	}
}

func TestSoundEquipmentAttributeSetAttributes(t *testing.T) {
	tests := []struct {
		name           string
		soundEquipAttr *SoundEquipmentAttribute
		expectedAttr   *SoundEquipmentAttribute
	}{
		{
			name:           "Test Case 1",
			soundEquipAttr: GetSoundEquipmentObject(1, 20000, []CodecType{CodecTypeAAC, CodecTypeLDAC}, 20),
			expectedAttr:   GetSoundEquipmentObject(1, 20000, []CodecType{CodecTypeAAC, CodecTypeLDAC}, 20),
		},
		{
			name:           "Test Case 2",
			soundEquipAttr: GetSoundEquipmentObject(1, 20000, []CodecType{CodecTypeSBC, CodecTypeAptx}, 29),
			expectedAttr:   GetSoundEquipmentObject(1, 20000, []CodecType{CodecTypeSBC, CodecTypeAptx}, 29),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.soundEquipAttr.SetAttributes().(*SoundEquipmentAttribute)
			assert.Equal(t, *tt.expectedAttr, *result, fmt.Sprintf("%s: Expected and actual SoundEquipmentAttributes are not equal", tt.name))
		})
	}
}

func GetWearablesObject(productID int, priceInCents uint64, strapTypes StrapType, stock int) *WearablesAttribute {
	return &WearablesAttribute{
		ProductID:    productID,
		PriceInCents: priceInCents,
		StrapTypes:   strapTypes,
		Stock:        stock,
	}
}

func TestWearablesAttributeSetAttributes(t *testing.T) {
	tests := []struct {
		name          string
		wearablesAttr *WearablesAttribute
		expectedAttr  *WearablesAttribute
	}{
		{
			name:          "Test Case 1",
			wearablesAttr: GetWearablesObject(1, 30000, StrapTypeSteel, 12),
			expectedAttr:  GetWearablesObject(1, 30000, StrapTypeSteel, 12),
		},
		{
			name:          "Test Case 2",
			wearablesAttr: GetWearablesObject(2, 29999, StrapTypeRubber, 30),
			expectedAttr:  GetWearablesObject(2, 29999, StrapTypeRubber, 30),
		},
		{
			name:          "Test Case 3",
			wearablesAttr: GetWearablesObject(3, 30000, StrapTypeLeather, 30),
			expectedAttr:  GetWearablesObject(3, 30000, StrapTypeLeather, 30),
		},
		{
			name:          "Test Case 4",
			wearablesAttr: GetWearablesObject(4, 29999, StrapTypeSilicone, 30),
			expectedAttr:  GetWearablesObject(4, 29999, StrapTypeSilicone, 30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.wearablesAttr.SetAttributes().(*WearablesAttribute)
			assert.Equal(t, *tt.expectedAttr, *result, fmt.Sprintf("%s: Expected and actual WearablesAttributes are not equal", tt.name))
		})
	}
}
