package service

import (
	"errors"
	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateNewCategoryRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.NewCategoryRequestDTO
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid Category - No Description",
			req:     domain.NewCategoryRequestDTO{Name: "Type-C", Description: ""},
			wantErr: false,
		},
		{
			name:    "Valid Category - With Description",
			req:     domain.NewCategoryRequestDTO{Name: "Sound Equipment", Description: "Great Quality"},
			wantErr: false,
		},
		{
			name:    "Invalid Category - Empty Name",
			req:     domain.NewCategoryRequestDTO{Name: "", Description: "Great Quality"},
			wantErr: true,
			errMsg:  "category name cannot be empty",
		},
		{
			name:    "Invalid Category - Invalid Characters in Name",
			req:     domain.NewCategoryRequestDTO{Name: "Sound@Equipment", Description: "Great Quality"},
			wantErr: true,
			errMsg:  "invalid characters in Category name field",
		},
		{
			name:    "Invalid Description - Description Too Long",
			req:     domain.NewCategoryRequestDTO{Name: "Phone", Description: string(make([]rune, 256))},
			wantErr: true,
			errMsg:  "category description must be less than 256 characters",
		},
		{
			name:    "Invalid Category and Description Too Long",
			req:     domain.NewCategoryRequestDTO{Name: "Sound@Equipment", Description: string(make([]rune, 256))},
			wantErr: true,
			errMsg:  "invalid characters in Category name field; category description must be less than 256 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNewCategoryRequest(tt.req)

			if tt.wantErr {
				var apiErr lib.APIError
				ok := errors.As(err, &apiErr)
				assert.True(t, ok, "Expected APIError type")
				assert.Equal(t, tt.errMsg, apiErr.AsMessage())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
