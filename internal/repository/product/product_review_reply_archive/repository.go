// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_archive_reply/product_archive_reply_repository.go
// 🧠 Concrete implementation of ProductReviewReplyArchiveRepoInterface

package product_review_reply_archive

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductReviewReplyArchiveRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewProductReviewReplyArchiveRepository(db *sql.DB) *ProductReviewReplyArchiveRepository {
	return &ProductReviewReplyArchiveRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductReviewReplyArchiveRepository) WithTx(
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

// 🧑 Validate seller
func (r *ProductReviewReplyArchiveRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Validate product ownership
func (r *ProductReviewReplyArchiveRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	productID uuid.UUID,
	sellerID uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
		ID:       productID,
		SellerID: sellerID,
	})
}

// 💬 Validate review
func (r *ProductReviewReplyArchiveRepository) GetProductReviewByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.ProductReview, error) {
	return r.q.GetProductReviewByID(ctx, id)
}

// 🔎 Validate existing reply
func (r *ProductReviewReplyArchiveRepository) GetReviewReplyByReviewIDAndSellerID(
	ctx context.Context,
	reviewID uuid.UUID,
	sellerID uuid.UUID,
) (sqlc.ProductReviewReply, error) {
	return r.q.GetReviewReplyByReviewIDAndSellerID(ctx, sqlc.GetReviewReplyByReviewIDAndSellerIDParams{
		ReviewID:     reviewID,
		SellerUserID: sellerID,
	})
}

// 🗑️ Archive reply
func (r *ProductReviewReplyArchiveRepository) ArchiveReviewReply(
	ctx context.Context,
	arg sqlc.ArchiveReviewReplyParams,
) error {
	return r.q.ArchiveReviewReply(ctx, arg)
}
