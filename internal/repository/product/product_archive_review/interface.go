package product_archive_review

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductArchiveReviewRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Validate customer
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Validate product
	GetProductByID(ctx context.Context, id uuid.UUID) (sqlc.Product, error)

	// 🔍 Validate review ownership + moderation
	GetProductReviewByIDAndProductIDAndReviewerID(
		ctx context.Context,
		reviewID uuid.UUID,
		productID uuid.UUID,
		reviewerID uuid.UUID,
	) (sqlc.ProductReview, error)

	// 🗃️ Archive review
	ArchiveProductReview(ctx context.Context, arg sqlc.ArchiveProductReviewParams) error
}
