package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ashtishad/ecommerce/lib"
	"log/slog"
)

type CategoryRepoDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewCategoryRepoDB(db *sql.DB, l *slog.Logger) *CategoryRepoDB {
	return &CategoryRepoDB{db, l}
}

func (d *CategoryRepoDB) CreateCategory(ctx context.Context, category Category) (*Category, lib.APIError) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		d.l.Error(lib.ErrTxBegin, "err", err)
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	defer func() {
		if err != nil {
			d.l.Error("unable to create category", "err", err.Error())
			if rbErr := tx.Rollback(); rbErr != nil {
				d.l.Warn("unable to rollback", "rollbackErr", rbErr)
			}
			return
		}
	}()

	if apiErr := d.checkCategoryNameExists(ctx, category.Name); apiErr != nil {
		return nil, apiErr
	}

	if err = tx.QueryRowContext(ctx,
		sqlInsertCategory, category.Name, category.Description).Scan(&category.CategoryID); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return d.findCategoryByID(ctx, category.CategoryID)
}

func (d *CategoryRepoDB) checkCategoryNameExists(ctx context.Context, categoryName string) lib.APIError {
	var existingCategoryName string

	err := d.db.QueryRowContext(ctx, sqlSelectCategoryName, categoryName).Scan(&existingCategoryName)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.l.Error("error checking existing category:", "err", err.Error())
		return lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if existingCategoryName != "" {
		d.l.Warn("category name already exists", "input", categoryName, "existing", existingCategoryName)
		return lib.NewDBFieldConflictError(fmt.Sprintf("category name already exists, input: %s existing:%s", categoryName, existingCategoryName))
	}

	return nil
}

// findCategoryByID takes categoryID and returns a single category record
// returns error(500 or 404) if internal server error happened.
func (d *CategoryRepoDB) findCategoryByID(ctx context.Context, categoryID int) (*Category, lib.APIError) {
	row := d.db.QueryRowContext(ctx, sqlSelectCategoryById, categoryID)

	var category Category
	err := row.Scan(&category.CategoryID,
		&category.CategoryUUID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			d.l.Error("category not found", "err", err.Error())
			return nil, lib.NewNotFoundError(lib.UnexpectedDatabaseErr)
		}
		d.l.Error(lib.ErrScanningRows, "err", err.Error())
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return &category, nil
}
