// ------------------------------------------------------------
// 📁 File: internal/services/product/edit_review_reply_service.go
// 🧠 Handles editing a seller's reply to a review.
//     - Validates seller
//     - Validates product ownership
//     - Validates review
//     - Validates existing reply
//     - Updates reply_text field

package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_review_reply_edit"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type EditReviewReplyInput struct {
	SellerUserID uuid.UUID
	ProductID    uuid.UUID
	ReviewID     uuid.UUID
	ReplyText    string
}

// ------------------------------------------------------------
// 📤 Response
type EditReviewReplyResult struct {
	ReviewID uuid.UUID
	Status   string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type EditReviewReplyServiceDeps struct {
	Repo repo.ProductReviewReplyEditRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service
type EditReviewReplyService struct {
	Deps EditReviewReplyServiceDeps
}

// 🚀 Constructor
func NewEditReviewReplyService(deps EditReviewReplyServiceDeps) *EditReviewReplyService {
	return &EditReviewReplyService{Deps: deps}
}

// 🚀 Entrypoint
func (s *EditReviewReplyService) Start(
	ctx context.Context,
	input EditReviewReplyInput,
) (*EditReviewReplyResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate seller
		user, err := q.GetUserByID(ctx, input.SellerUserID)
		if err != nil {
			return errors.NewNotFoundError("seller")
		}
		if user.IsBanned || user.IsArchived || !user.IsSellerProfileApproved || !user.IsSellerProfileCreated {
			return errors.NewAuthError("unauthorized seller")
		}

		// ------------------------------------------------------------
		// Step 2: Validate product ownership
		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.SellerUserID,
		})
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsBanned || product.IsArchived || !product.IsApproved {
			return errors.NewValidationError("product", "not approved, banned, or archived")
		}

		// ------------------------------------------------------------
		// Step 3: Validate product review
		review, err := q.GetProductReviewByID(ctx, input.ReviewID)
		if err != nil {
			return errors.NewNotFoundError("review")
		}
		if review.IsBanned || review.IsArchived {
			return errors.NewValidationError("review", "banned or archived")
		}

		// ------------------------------------------------------------
		// Step 4: Validate existing reply
		reply, err := q.GetReviewReplyByReviewIDAndSellerID(ctx, sqlc.GetReviewReplyByReviewIDAndSellerIDParams{
			ReviewID:     input.ReviewID,
			SellerUserID: input.SellerUserID,
		})
		if err != nil {
			return errors.NewNotFoundError("review_reply")
		}
		if reply.IsBanned || reply.IsArchived {
			return errors.NewValidationError("reply", "banned or archived")
		}

		// ------------------------------------------------------------
		// Step 5: Update reply
		err = q.UpdateReviewReplyText(ctx, sqlc.UpdateReviewReplyTextParams{
			ReviewID:     input.ReviewID,
			SellerUserID: input.SellerUserID,
			ReplyText:    input.ReplyText,
			IsEdited:     true,
			UpdatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("product_review_replies.update", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &EditReviewReplyResult{
		ReviewID: input.ReviewID,
		Status:   "reply_updated",
	}, nil
}
