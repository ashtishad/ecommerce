package service

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/ashtishad/ecommerce/lib"
	"github.com/ashtishad/ecommerce/product-api/internal/domain"
)

// ValidateNewCategoryRequest validates the new category request data.
// It returns a combined error message if any of the validations fail.
//
// - Name must not be empty and should only contain alphanumeric characters, spaces, or certain symbols like '-' and '_'.
// - Description must be less than 256 characters in length.
//
// If any validation rule is not met, it returns an APIError with all error messages combined.
func ValidateNewCategoryRequest(req domain.NewCategoryRequestDTO) lib.APIError {
	var errMessages []string

	if req.Name == "" {
		errMessages = append(errMessages, "category name cannot be empty")
	}

	nameRegex := `^[A-Za-z0-9\s\-_&]*$`
	if !regexp.MustCompile(nameRegex).MatchString(req.Name) {
		errMessages = append(errMessages, "invalid characters in Category name field")
	}

	if utf8.RuneCountInString(req.Description) > 255 {
		errMessages = append(errMessages, "category description must be less than 256 characters")
	}

	if len(errMessages) > 0 {
		return lib.NewBadRequestError(strings.Join(errMessages, "; "))
	}

	return nil
}
