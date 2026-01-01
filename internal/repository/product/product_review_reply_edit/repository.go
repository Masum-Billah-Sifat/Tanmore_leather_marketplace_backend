// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_edit_reply/product_edit_reply_repository.go
// 🧠 Concrete implementation of ProductEditReplyRepoInterface.

package product_edit_reply

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductReviewReplyEditRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewProductReviewReplyEditRepository(db *sql.DB) *ProductReviewReplyEditRepository {
	return &ProductReviewReplyEditRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductReviewReplyEditRepository) WithTx(
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
func (r *ProductReviewReplyEditRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Validate product ownership
func (r *ProductReviewReplyEditRepository) GetProductByIDAndSellerID(
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
func (r *ProductReviewReplyEditRepository) GetProductReviewByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.ProductReview, error) {
	return r.q.GetProductReviewByID(ctx, id)
}

// 🔎 Validate existing reply
func (r *ProductReviewReplyEditRepository) GetReviewReplyByReviewIDAndSellerID(
	ctx context.Context,
	reviewID uuid.UUID,
	sellerID uuid.UUID,
) (sqlc.ProductReviewReply, error) {
	return r.q.GetReviewReplyByReviewIDAndSellerID(ctx, sqlc.GetReviewReplyByReviewIDAndSellerIDParams{
		ReviewID:     reviewID,
		SellerUserID: sellerID,
	})
}

// ✏️ Update reply
func (r *ProductReviewReplyEditRepository) UpdateReviewReplyText(
	ctx context.Context,
	arg sqlc.UpdateReviewReplyTextParams,
) error {
	return r.q.UpdateReviewReplyText(ctx, arg)
}
