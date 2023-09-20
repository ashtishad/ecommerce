package domain

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/ashtishad/ecommerce/lib"
)

type Brand struct {
	BrandID   int64
	BrandUUID string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BrandResponseDTO struct {
	BrandUUID string    `json:"brandUuid"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Brand) ToBrandResponseDTO() *BrandResponseDTO {
	return &BrandResponseDTO{
		BrandUUID: b.BrandUUID,
		Name:      b.Name,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

type BrandRepository interface {
	GetAllBrands(ctx context.Context, status string) ([]Brand, lib.APIError)
}

type BrandRepoDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewBrandRepoDB(db *sql.DB, l *slog.Logger) *BrandRepoDB {
	return &BrandRepoDB{
		db: db,
		l:  l,
	}
}

func (d *BrandRepoDB) GetAllBrands(ctx context.Context, status string) ([]Brand, lib.APIError) {
	query := `SELECT brand_id,brand_uuid,name,status, created_at, updated_at FROM brands where status = $1`
	rows, err := d.db.QueryContext(ctx, query, status)

	if err != nil {
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	defer rows.Close()

	res := make([]Brand, 0)

	for rows.Next() {
		var b Brand

		err := rows.Scan(&b.BrandID, &b.BrandUUID, &b.Name, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, lib.NewNotFoundError("brands not found")
			}

			d.l.Error("unable to scan rows", "err", err)

			return nil, lib.NewInternalServerError(lib.ErrScanningRows, err)
		}

		res = append(res, b)
	}

	if rErr := rows.Close(); rErr != nil {
		d.l.Error("unable to close rows", "err", rErr)
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, rErr)
	}

	if err := rows.Err(); err != nil {
		d.l.Error("unexpected error", "err", err)
		return nil, lib.NewInternalServerError(lib.UnexpectedDatabaseErr, err)
	}

	return res, nil
}
