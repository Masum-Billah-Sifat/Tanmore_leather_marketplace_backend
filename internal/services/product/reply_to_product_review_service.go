// ------------------------------------------------------------
// 📁 File: internal/services/product/reply_to_product_review_service.go
// 🧠 Handles replying to a review for a product.
//     - Validates seller
//     - Validates product ownership
//     - Validates target review
//     - Ensures no existing reply
//     - Inserts reply into product_review_replies
//     - Returns reply_id and review_id

package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_reply_review"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type ReplyToReviewInput struct {
	SellerID      uuid.UUID
	ProductID     uuid.UUID
	ReviewID      uuid.UUID
	ReplyText     string
	ReplyImageURL string // optional
}

// ------------------------------------------------------------
// 📤 Result to return
type ReplyToReviewResult struct {
	ReviewID uuid.UUID
	ReplyID  uuid.UUID
	Status   string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type ReplyToReviewServiceDeps struct {
	Repo repo.ProductReplyReviewRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type ReplyToReviewService struct {
	Deps ReplyToReviewServiceDeps
}

// 🚀 Constructor
func NewReplyToReviewService(deps ReplyToReviewServiceDeps) *ReplyToReviewService {
	return &ReplyToReviewService{Deps: deps}
}

// 🚀 Entrypoint
func (s *ReplyToReviewService) Start(
	ctx context.Context,
	input ReplyToReviewInput,
) (*ReplyToReviewResult, error) {

	now := timeutil.NowUTC()
	var replyID uuid.UUID

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate seller
		user, err := q.GetUserByID(ctx, input.SellerID)
		if err != nil {
			return errors.NewNotFoundError("seller")
		}
		if user.IsArchived {
			return errors.NewAuthError("seller is archived")
		}
		if user.IsBanned {
			return errors.NewAuthError("seller is banned")
		}
		if !user.IsSellerProfileCreated || !user.IsSellerProfileApproved {
			return errors.NewAuthError("seller profile not ready or not approved")
		}

		// ------------------------------------------------------------
		// Step 2: Validate product ownership
		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.SellerID,
		})
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsArchived || product.IsBanned || !product.IsApproved {
			return errors.NewValidationError("product", "not approved, banned, or archived")
		}

		// ------------------------------------------------------------
		// Step 3: Validate review
		review, err := q.GetProductReviewByID(ctx, input.ReviewID)
		if err != nil {
			return errors.NewNotFoundError("review")
		}
		if review.IsArchived || review.IsBanned {
			return errors.NewValidationError("review", "archived or banned")
		}

		// ------------------------------------------------------------
		// Step 4: Check existing reply
		existingReply, err := q.GetReviewReplyByReviewID(ctx, input.ReviewID)
		if err == nil {
			if existingReply.IsArchived || existingReply.IsBanned {
				return errors.NewValidationError("reply", "already exists but archived or banned")
			}
			return errors.NewServerError("a reply already exists for this review")
		}

		// ------------------------------------------------------------
		// Step 5: Insert reply
		replyUUID := uuidutil.New()
		replyID = replyUUID

		_, err = q.InsertReviewReply(ctx, sqlc.InsertReviewReplyParams{
			ID:            replyUUID,
			ReviewID:      input.ReviewID,
			SellerUserID:  input.SellerID,
			ReplyText:     input.ReplyText,
			ReplyImageUrl: sqlnull.String(input.ReplyImageURL),
			IsEdited:      false,
			IsArchived:    false,
			IsBanned:      false,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
		if err != nil {
			return errors.NewTableError("product_review_replies.insert", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ReplyToReviewResult{
		ReviewID: input.ReviewID,
		ReplyID:  replyID,
		Status:   "reply_submitted",
	}, nil
}
