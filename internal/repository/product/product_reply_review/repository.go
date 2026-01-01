// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_reply_review/product_reply_review_repository.go
// 🧠 Concrete implementation of ProductReplyReviewRepoInterface

package product_reply_review

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductReplyReviewRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewProductReplyReviewRepository(db *sql.DB) *ProductReplyReviewRepository {
	return &ProductReplyReviewRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *ProductReplyReviewRepository) WithTx(
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

// ✅ Seller validation
func (r *ProductReplyReviewRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📦 Product ownership validation
func (r *ProductReplyReviewRepository) GetProductByIDAndSellerID(
	ctx context.Context,
	productID uuid.UUID,
	sellerID uuid.UUID,
) (sqlc.Product, error) {
	return r.q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
		ID:       productID,
		SellerID: sellerID,
	})
}

// 🔍 Review validation
func (r *ProductReplyReviewRepository) GetProductReviewByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.ProductReview, error) {
	return r.q.GetProductReviewByID(ctx, id)
}

// 🔄 Check existing reply
func (r *ProductReplyReviewRepository) GetReviewReplyByReviewID(
	ctx context.Context,
	reviewID uuid.UUID,
) (sqlc.ProductReviewReply, error) {
	return r.q.GetReviewReplyByReviewID(ctx, reviewID)
}

// ✍️ Insert reply
func (r *ProductReplyReviewRepository) InsertReviewReply(
	ctx context.Context,
	arg sqlc.InsertReviewReplyParams,
) (uuid.UUID, error) {
	return r.q.InsertReviewReply(ctx, arg)
}
