package domain

import (
	"context"

	"github.com/ashtishad/ecommerce/lib"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category Category) (*Category, lib.APIError)
	CreateSubCategory(ctx context.Context, subCategory Category, parentCategoryUUID string) (*Category, lib.APIError)
	GetAllCategoriesWithHierarchy(ctx context.Context) ([]*Category, lib.APIError)

	checkCategoryNameExists(ctx context.Context, categoryName string) lib.APIError
	findCategoryByID(ctx context.Context, categoryID int) (*Category, lib.APIError)
}
