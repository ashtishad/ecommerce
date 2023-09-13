package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ashtishad/ecommerce/lib"
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

	defer rollBackOnError(tx, d.l, &err)

	if apiErr := d.checkCategoryNameExists(ctx, category.Name); apiErr != nil {
		return nil, apiErr
	}

	if err = d.executeInsertCategory(ctx, tx, &category); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return d.findCategoryByID(ctx, category.CategoryID)
}

// ExecuteInsertCategory executes the SQL query to insert a new category
// Scans category id and assigns to category domain object
func (d *CategoryRepoDB) executeInsertCategory(ctx context.Context, tx *sql.Tx, category *Category) error {
	var categoryID int
	err := tx.QueryRowContext(
		ctx,
		sqlInsertCategory,
		category.Name,
		category.Description,
	).Scan(&categoryID)

	if err != nil {
		return fmt.Errorf("unable to execute insert category: %w", err)
	}

	category.CategoryID = categoryID

	return nil
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
	row := d.db.QueryRowContext(ctx, sqlSelectCategoryByID, categoryID)

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

func (d *CategoryRepoDB) CreateSubCategory(ctx context.Context, subCategory Category, parentCategoryUUID string) (*Category, lib.APIError) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		d.l.Error(lib.ErrTxBegin, "err", err)
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	defer rollBackOnError(tx, d.l, &err)

	// validate parentCategoryUUID and get its ID
	parentCategoryID, apiErr := d.validateAndGetParentID(ctx, tx, parentCategoryUUID, subCategory.Name)
	if apiErr != nil {
		return nil, apiErr
	}

	// first insert sub-category
	if err = d.executeInsertCategory(ctx, tx, &subCategory); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	// calculate parent and insert category relationships
	if err = d.insertCategoryRelationship(ctx, tx, parentCategoryID, subCategory.CategoryID); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return d.findCategoryByID(ctx, subCategory.CategoryID)
}

func rollBackOnError(tx *sql.Tx, l *slog.Logger, err *error) {
	if *err != nil {
		l.Error("unable to complete operation", "err", (*err).Error())

		if rbErr := tx.Rollback(); rbErr != nil {
			l.Warn("unable to rollback", "rollbackErr", rbErr)
		}
	}
}

// validateAndGetParentID validates the given parent category UUID, checks for category name uniqueness,
// and returns its ID, otherwise, it returns error 404,500 or 409.
func (d *CategoryRepoDB) validateAndGetParentID(ctx context.Context, tx *sql.Tx, parentCategoryUUID string, categoryName string) (int, lib.APIError) {
	var (
		parentCategoryID int
		nameExists       bool
	)

	err := tx.QueryRowContext(ctx,
		sqlValidateUUIDGetCatID,
		parentCategoryUUID, categoryName).Scan(&parentCategoryID, &nameExists)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, lib.NewNotFoundError("parent category not found")
	} else if err != nil {
		d.l.Error("database error:", "err", err)
		return 0, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	if nameExists {
		return 0, lib.NewDBFieldConflictError("category name already exists")
	}

	return parentCategoryID, nil
}

// insertCategoryRelationship inserts a new row into the category_relationships table and calculates the level.
func (d *CategoryRepoDB) insertCategoryRelationship(ctx context.Context, tx *sql.Tx, parentCategoryID, subCategoryID int) lib.APIError {
	_, err := tx.ExecContext(ctx, sqlInsertWithLevelCalculation, parentCategoryID, parentCategoryID, subCategoryID)
	if err != nil {
		d.l.Error("failed to insert into category_relationships:", "err", err)
		return lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return nil
}

func (d *CategoryRepoDB) GetAllCategoriesWithHierarchy(ctx context.Context) ([]*Category, lib.APIError) {
	rows, err := d.db.QueryContext(ctx, sqlGetAllCategoriesWithHierarchy)
	if err != nil {
		d.l.Error("failed to query get all categories:", "err", err)
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	defer closeRows(rows, d.l)

	return d.BuildTree(rows)
}

func (d *CategoryRepoDB) BuildTree(rows *sql.Rows) ([]*Category, lib.APIError) {
	var categories []*Category

	for rows.Next() {
		var category Category
		err := rows.Scan(&category.CategoryUUID, &category.ParentCategoryUUID, &category.Level, &category.Name)

		if err != nil {
			d.l.Error("failed to scan rows:", "err", err)
			return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
		}

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		d.l.Error("unexpected error on BuildTree", "err", err)
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	emptyParent := sql.NullString{Valid: false}

	return buildTree(categories, emptyParent), nil
}

func buildTree(categories []*Category, parentUUID sql.NullString) []*Category {
	var tree []*Category

	for _, c := range categories {
		if parentUUID.Valid && c.ParentCategoryUUID.Valid && parentUUID.String == c.ParentCategoryUUID.String {
			children := buildTree(categories, sql.NullString{String: c.CategoryUUID, Valid: true})
			if len(children) > 0 {
				sub := make([]Category, len(children))
				for i := range children {
					sub[i] = *children[i]
				}

				c.Subcategories = sub
			}

			tree = append(tree, c)
		} else if !parentUUID.Valid && !c.ParentCategoryUUID.Valid {
			children := buildTree(categories, sql.NullString{String: c.CategoryUUID, Valid: true})
			if len(children) > 0 {
				sub := make([]Category, len(children))
				for i := range children {
					sub[i] = *children[i]
				}
				c.Subcategories = sub
			}
			tree = append(tree, c)
		}
	}

	return tree
}

func closeRows(rows *sql.Rows, l *slog.Logger) {
	if rcErr := rows.Close(); rcErr != nil {
		l.Warn("error closing rows", "err", rcErr)
	}
}
