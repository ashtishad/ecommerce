package domain

import (
	"context"

	"github.com/ashtishad/ecommerce/lib"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, categoryUUID string, p Product) (*Product, lib.APIError)
}
