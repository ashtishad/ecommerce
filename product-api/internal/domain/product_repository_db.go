package domain

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ashtishad/ecommerce/lib"
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

// CreateProduct takes Product and its category_uuid, then find its category_id and root category_id of level 0
// inserts product if everything goes alright
// return 404 or 500 error
func (d ProductRepoDB) CreateProduct(ctx context.Context, categoryUUID string, p Product) (*Product, lib.APIError) {
	var catID, rootCatID string
	if err := d.db.QueryRowContext(ctx, sqlGetCategoryIDAndRootCategoryID, categoryUUID).Scan(&catID, &rootCatID); err != nil {
		d.l.Error("unable to get category id and root id", "err", err.Error())
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, lib.NewInternalServerError(lib.ErrTxBegin, err)
	}

	defer rollBackOnError(tx, d.l, &err)

	if err := d.db.QueryRowContext(ctx, sqlInsertProductRetID, p.Name, catID, rootCatID).Scan(&p.ProductID); err != nil {
		d.l.Error("unable to insert product id", "err", err.Error())
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, lib.NewInternalServerError("unable to commit", err)
	}

	return d.findProductByID(ctx, p.ProductID)
}

func (d ProductRepoDB) findProductByID(ctx context.Context, id int) (*Product, lib.APIError) {
	row := d.db.QueryRowContext(ctx, sqlFindProductByID, id)

	var p Product
	err := row.Scan(&p.ProductUUID,
		&p.Name, &p.CategoryID, &p.RootCategoryID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d.l.Error("product not found", "id", id, "err", err.Error())
			return nil, lib.NewNotFoundError("product not found")
		}

		d.l.Error("error scanning rows", "err", err.Error())

		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return &p, nil
}
