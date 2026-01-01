// ------------------------------------------------------------
// 📁 File: internal/services/product/edit_product_review_service.go
// 🧠 Handles editing a review for a product.
//     - Validates customer
//     - Validates product status
//     - Validates review ownership and moderation
//     - Updates review text
//     - Returns review_id and product_id

package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_edit_review"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type EditProductReviewInput struct {
	UserID     uuid.UUID
	ProductID  uuid.UUID
	ReviewID   uuid.UUID
	ReviewText string
}

// ------------------------------------------------------------
// 📤 Result to return
type EditProductReviewResult struct {
	ProductID uuid.UUID
	ReviewID  uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type EditProductReviewServiceDeps struct {
	Repo repo.ProductEditReviewRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type EditProductReviewService struct {
	Deps EditProductReviewServiceDeps
}

// 🚀 Constructor
func NewEditProductReviewService(deps EditProductReviewServiceDeps) *EditProductReviewService {
	return &EditProductReviewService{Deps: deps}
}

// 🚀 Entrypoint
func (s *EditProductReviewService) Start(
	ctx context.Context,
	input EditProductReviewInput,
) (*EditProductReviewResult, error) {

	now := timeutil.NowUTC()

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
		// Step 3: Validate review ownership & moderation
		_, err = q.GetProductReviewByIDAndProductIDAndReviewerID(ctx, sqlc.GetProductReviewByIDAndProductIDAndReviewerIDParams{
			ID:             input.ReviewID,
			ProductID:      input.ProductID,
			ReviewerUserID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("review")
		}
		// You can later expand to check `is_archived` / `is_banned` if needed here
		// or include them in the SQL query's WHERE clause.

		// ------------------------------------------------------------
		// Step 4: Update review text
		err = q.UpdateProductReviewText(ctx, sqlc.UpdateProductReviewTextParams{
			ID:             input.ReviewID,
			ProductID:      input.ProductID,
			ReviewerUserID: input.UserID,
			ReviewText:     input.ReviewText,
			IsEdited:       true,
			UpdatedAt:      now,
		})
		if err != nil {
			return errors.NewTableError("product_reviews.update", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &EditProductReviewResult{
		ProductID: input.ProductID,
		ReviewID:  input.ReviewID,
		Status:    "review_updated",
	}, nil
}
