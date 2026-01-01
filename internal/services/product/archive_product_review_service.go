package product

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_archive_review"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type ArchiveProductReviewInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	ReviewID  uuid.UUID
}

// ------------------------------------------------------------
// 📤 Output to return
type ArchiveProductReviewResult struct {
	ProductID uuid.UUID
	ReviewID  uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type ArchiveProductReviewServiceDeps struct {
	Repo repo.ProductArchiveReviewRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type ArchiveProductReviewService struct {
	Deps ArchiveProductReviewServiceDeps
}

// 🚀 Constructor
func NewArchiveProductReviewService(deps ArchiveProductReviewServiceDeps) *ArchiveProductReviewService {
	return &ArchiveProductReviewService{Deps: deps}
}

// 🚀 Entrypoint
func (s *ArchiveProductReviewService) Start(
	ctx context.Context,
	input ArchiveProductReviewInput,
) (*ArchiveProductReviewResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Validate Customer
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

		// Step 2: Validate Product
		product, err := q.GetProductByID(ctx, input.ProductID)
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsArchived || product.IsBanned || !product.IsApproved {
			return errors.NewValidationError("product", "not approved, banned, or archived")
		}

		// Step 3: Validate Review
		review, err := q.GetProductReviewByIDAndProductIDAndReviewerID(ctx, sqlc.GetProductReviewByIDAndProductIDAndReviewerIDParams{
			ID:             input.ReviewID,
			ProductID:      input.ProductID,
			ReviewerUserID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("review")
		}
		if review.IsArchived {
			return errors.NewValidationError("review", "review is already archived")
		}
		if review.IsBanned {
			return errors.NewValidationError("review", "review is banned")
		}

		// Step 4: Archive the review
		err = q.ArchiveProductReview(ctx, sqlc.ArchiveProductReviewParams{
			ID:         input.ReviewID,
			IsArchived: true,
			UpdatedAt:  now,
		})
		if err != nil {
			return errors.NewTableError("product_reviews.archive", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ArchiveProductReviewResult{
		ProductID: input.ProductID,
		ReviewID:  input.ReviewID,
		Status:    "review_archived",
	}, nil
}
