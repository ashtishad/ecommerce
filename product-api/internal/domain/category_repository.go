package domain

import (
	"context"
	"github.com/ashtishad/ecommerce/lib"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category Category) (*Category, lib.APIError)

	findCategoryByID(ctx context.Context, categoryID int) (*Category, lib.APIError)
}
