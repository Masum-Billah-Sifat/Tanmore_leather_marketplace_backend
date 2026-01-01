// ------------------------------------------------------------
// 📁 File: internal/repository/product/fetch_by_category/fetch_by_category_repository.go
// 🧠 Concrete implementation of FetchByCategoryRepoInterface
//     Uses SQLC for category check, recursive leaf lookup, and variant index filtering.

package fetch_by_category

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 FetchByCategoryRepository implements FetchByCategoryRepoInterface
type FetchByCategoryRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewFetchByCategoryRepository(db *sql.DB) *FetchByCategoryRepository {
	return &FetchByCategoryRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *FetchByCategoryRepository) WithTx(
	ctx context.Context,
	fn func(q *sqlc.Queries) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := sqlc.New(tx)

	if err := fn(qtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// 🔍 Category validation
func (r *FetchByCategoryRepository) GetCategoryByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.Category, error) {
	return r.q.GetCategoryByID(ctx, id)
}

// 🌿 Recursive leaf category lookup
func (r *FetchByCategoryRepository) GetAllLeafCategoryIDsByRoot(
	ctx context.Context,
	rootID uuid.UUID,
) ([]uuid.UUID, error) {
	return r.q.GetAllLeafCategoryIDsByRoot(ctx, rootID)
}

// 🛍️ Fetch valid product variant index rows
func (r *FetchByCategoryRepository) GetProductVariantIndexesByCategoryIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]sqlc.GetProductVariantIndexesByCategoryIDsRow, error) {
	return r.q.GetProductVariantIndexesByCategoryIDs(ctx, ids)
}
