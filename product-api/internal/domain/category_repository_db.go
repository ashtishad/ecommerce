package domain

import (
	"context"
	"database/sql"
	"errors"
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
	var existingCategoryName string

	err := d.db.QueryRowContext(ctx,
		"SELECT category_id FROM categories WHERE LOWER(name) = LOWER($1)",
		category.Name).Scan(&existingCategoryName)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.l.Error("error checking existing category:", "err", err.Error())
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if existingCategoryName != "" {
		d.l.Warn("category name already exists", "input_cat_name", category.Name, "exist_cat_name", existingCategoryName)
		return nil, lib.NewDBFieldConflictError("category name already exists", category.Name, existingCategoryName)
	}

	err = d.db.QueryRowContext(ctx,
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id",
		category.Name, category.Description).Scan(&category.CategoryID)

	if err != nil {
		d.l.Error("error inserting new category:", "err", err.Error())
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return d.findCategoryByID(ctx, category.CategoryID)
}

// findCategoryByID takes categoryID and returns a single category record
// returns error(500 or 404) if internal server error happened.
func (d *CategoryRepoDB) findCategoryByID(ctx context.Context, categoryID int) (*Category, lib.APIError) {
	query := `SELECT category_id,category_uuid,name, description,status,created_at,updated_at FROM categories where category_id= $1`
	row := d.db.QueryRowContext(ctx, query, categoryID)

	var category Category
	err := row.Scan(&category.CategoryID,
		&category.CategoryUUID,
		&category.Name,
		&category.Description,
		&category.Status,
		&category.CreatedAt,
		&category.UpdatedAt)
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				d.l.Error("category not found", "err", err.Error())
				return nil, lib.NewNotFoundError(lib.UnexpectedDatabaseErr)
			}
			d.l.Error(lib.ErrScanningRows, "err", err.Error())
			return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
		}
	}

	return &category, nil
}
