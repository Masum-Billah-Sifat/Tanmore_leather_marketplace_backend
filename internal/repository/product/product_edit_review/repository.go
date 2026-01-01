// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_edit_review/product_edit_review_repository.go
// 🧠 Concrete implementation of ProductEditReviewRepoInterface.
//     Uses SQLC to validate user/product/review and update review.

package product_edit_review

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductEditReviewRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductEditReviewRepository(db *sql.DB) *ProductEditReviewRepository {
	return &ProductEditReviewRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductEditReviewRepository) WithTx(
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

// 👤 Validate customer
func (r *ProductEditReviewRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Validate product
func (r *ProductEditReviewRepository) GetProductByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByID(ctx, id)
}

// 🔍 Validate review ownership
func (r *ProductEditReviewRepository) GetProductReviewByIDAndProductIDAndReviewerID(
	ctx context.Context,
	arg sqlc.GetProductReviewByIDAndProductIDAndReviewerIDParams,
) (sqlc.ProductReview, error) {
	return r.q.GetProductReviewByIDAndProductIDAndReviewerID(ctx, arg)
}

// ✏️ Update review text
func (r *ProductEditReviewRepository) UpdateProductReviewText(
	ctx context.Context,
	arg sqlc.UpdateProductReviewTextParams,
) error {
	return r.q.UpdateProductReviewText(ctx, arg)
}
