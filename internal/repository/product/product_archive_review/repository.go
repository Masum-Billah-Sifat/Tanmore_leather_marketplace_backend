package product_archive_review

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductArchiveReviewRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewProductArchiveReviewRepository(db *sql.DB) *ProductArchiveReviewRepository {
	return &ProductArchiveReviewRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductArchiveReviewRepository) WithTx(
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

// 👤 Get customer
func (r *ProductArchiveReviewRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Get product
func (r *ProductArchiveReviewRepository) GetProductByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByID(ctx, id)
}

// 🔍 Validate review ownership + moderation
func (r *ProductArchiveReviewRepository) GetProductReviewByIDAndProductIDAndReviewerID(
	ctx context.Context,
	reviewID uuid.UUID,
	productID uuid.UUID,
	reviewerID uuid.UUID,
) (sqlc.ProductReview, error) {
	return r.q.GetProductReviewByIDAndProductIDAndReviewerID(ctx, sqlc.GetProductReviewByIDAndProductIDAndReviewerIDParams{
		ID:             reviewID,
		ProductID:      productID,
		ReviewerUserID: reviewerID,
	})
}

// 🗃️ Archive review
func (r *ProductArchiveReviewRepository) ArchiveProductReview(
	ctx context.Context,
	arg sqlc.ArchiveProductReviewParams,
) error {
	return r.q.ArchiveProductReview(ctx, arg)
}
