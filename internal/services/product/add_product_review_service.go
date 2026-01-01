// ------------------------------------------------------------
// 📁 File: internal/services/product/add_product_review_service.go
// 🧠 Handles submitting a review for a product.
//     - Validates customer
//     - Validates product status
//     - Inserts review row into product_reviews
//     - Returns review_id and product_id

package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_add_review"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type AddProductReviewInput struct {
	UserID         uuid.UUID
	ProductID      uuid.UUID
	ReviewText     string
	ReviewImageURL string // optional
}

// ------------------------------------------------------------
// 📤 Result to return
type AddProductReviewResult struct {
	ProductID uuid.UUID
	ReviewID  uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type AddProductReviewServiceDeps struct {
	Repo repo.ProductAddReviewRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type AddProductReviewService struct {
	Deps AddProductReviewServiceDeps
}

// 🚀 Constructor
func NewAddProductReviewService(deps AddProductReviewServiceDeps) *AddProductReviewService {
	return &AddProductReviewService{Deps: deps}
}

// 🚀 Entrypoint
func (s *AddProductReviewService) Start(
	ctx context.Context,
	input AddProductReviewInput,
) (*AddProductReviewResult, error) {

	now := timeutil.NowUTC()
	var reviewID uuid.UUID

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate customer
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived {
			return errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			return errors.NewAuthError("user is banned")
		}

		// ------------------------------------------------------------
		// Step 2: Validate product
		product, err := q.GetProductByID(ctx, input.ProductID)
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsArchived || product.IsBanned || !product.IsApproved {
			return errors.NewValidationError("product", "not approved, banned, or archived")
		}

		// ------------------------------------------------------------
		// Step 3: Insert review
		reviewUUID := uuidutil.New()
		reviewID = reviewUUID

		_, err = q.InsertProductReview(ctx, sqlc.InsertProductReviewParams{
			ID:             reviewUUID,
			ProductID:      input.ProductID,
			ReviewerUserID: input.UserID,
			ReviewText:     input.ReviewText,
			ReviewImageUrl: sqlnull.String(input.ReviewImageURL),
			IsEdited:       false,
			IsArchived:     false,
			IsBanned:       false,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if err != nil {
			return errors.NewTableError("product_reviews.insert", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AddProductReviewResult{
		ProductID: input.ProductID,
		ReviewID:  reviewID,
		Status:    "review_submitted",
	}, nil
}
