// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_add_review/product_add_review_repository.go
// 🧠 Concrete implementation of ProductAddReviewRepoInterface.
//     Uses SQLC to validate customer, product, and insert review.

package product_add_review

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 ProductAddReviewRepository implements ProductAddReviewRepoInterface
type ProductAddReviewRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductAddReviewRepository(db *sql.DB) *ProductAddReviewRepository {
	return &ProductAddReviewRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductAddReviewRepository) WithTx(
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

// 👤 Fetch customer by ID
func (r *ProductAddReviewRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Fetch product by ID
func (r *ProductAddReviewRepository) GetProductByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByID(ctx, id)
}

// ✍️ Insert product review
func (r *ProductAddReviewRepository) InsertProductReview(
	ctx context.Context,
	arg sqlc.InsertProductReviewParams,
) (uuid.UUID, error) {
	return r.q.InsertProductReview(ctx, arg)
}
