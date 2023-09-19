package domain

import (
	"database/sql"
	"log/slog"
)

type ProductRepoDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewProductRepoDB(db *sql.DB, l *slog.Logger) *ProductRepoDB {
	return &ProductRepoDB{
		db: db,
		l:  l,
	}
}

// func (d *ProductRepoDB) CreateProduct(ctx context.Context, p Product, catUUID string) (*Product, lib.APIError) {
// 	// retrieve category id from uuid
//
// 	// insert product and it's attributes
//
// 	// get product from id
//
// 	return nil, nil
// }
//
// func (d *ProductRepoDB) findProductByID(ctx context.Context, productID int) (*Product, lib.APIError) {
// 	query := `SELECT `
// }
