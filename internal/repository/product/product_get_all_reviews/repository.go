// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_get_all_reviews/product_get_all_reviews_repository.go
// 🧠 Concrete implementation of ProductGetAllReviewsRepoInterface

package product_get_all_reviews

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📆 Repository struct
type ProductGetAllReviewsRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductGetAllReviewsRepository(db *sql.DB) *ProductGetAllReviewsRepository {
	return &ProductGetAllReviewsRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper (optional use)
func (r *ProductGetAllReviewsRepository) WithTx(
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

// 🛄 Product validation
func (r *ProductGetAllReviewsRepository) GetProductByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByID(ctx, id)
}

func (r *ProductGetAllReviewsRepository) GetAllReviewsByProductID(
	ctx context.Context,
	arg sqlc.GetAllReviewsByProductIDParams,
) ([]sqlc.GetAllReviewsByProductIDRow, error) {
	return r.q.GetAllReviewsByProductID(ctx, arg)
}

func (r *ProductGetAllReviewsRepository) GetRepliesByReviewIDs(
	ctx context.Context,
	reviewIDs []uuid.UUID,
) ([]sqlc.GetRepliesByReviewIDsRow, error) {
	return r.q.GetRepliesByReviewIDs(ctx, reviewIDs)
}
