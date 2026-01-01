// ------------------------------------------------------------
// 📁 File: internal/services/product/get_all_product_reviews_service.go
// 🧠 Handles fetching all reviews for a product (public endpoint)
//     - Validates product
//     - Fetches all non-banned, non-archived reviews
//     - Batch fetches replies for those reviews
//     - Merges replies into each review response

package product

import (
	"context"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_get_all_reviews"

	"tanmore_backend/pkg/errors"

	"tanmore_backend/pkg/sqlnull"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler

type GetAllProductReviewsInput struct {
	ProductID uuid.UUID
	Page      int
	Limit     int
}

// ------------------------------------------------------------
// 📤 Response Structs

type ReviewWithReply struct {
	ReviewID       uuid.UUID      `json:"review_id"`
	ReviewerUserID uuid.UUID      `json:"reviewer_user_id"`
	ReviewText     string         `json:"review_text"`
	ReviewImageURL *string        `json:"review_image_url,omitempty"`
	IsEdited       bool           `json:"is_edited"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Reply          *ReplyResponse `json:"reply,omitempty"`
}

type ReplyResponse struct {
	ReplyID       uuid.UUID `json:"reply_id"`
	SellerUserID  uuid.UUID `json:"seller_user_id"`
	ReplyText     string    `json:"reply_text"`
	ReplyImageURL *string   `json:"reply_image_url,omitempty"`
	IsEdited      bool      `json:"is_edited"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ------------------------------------------------------------
// 📤 Final Output

type GetAllProductReviewsResult struct {
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalItems int               `json:"total_items"`
	Items      []ReviewWithReply `json:"items"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type GetAllProductReviewsServiceDeps struct {
	Repo repo.ProductGetAllReviewsRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service

type GetAllProductReviewsService struct {
	Deps GetAllProductReviewsServiceDeps
}

// 🚀 Constructor
func NewGetAllProductReviewsService(deps GetAllProductReviewsServiceDeps) *GetAllProductReviewsService {
	return &GetAllProductReviewsService{Deps: deps}
}

// 🚀 Entrypoint
func (s *GetAllProductReviewsService) Start(
	ctx context.Context,
	input GetAllProductReviewsInput,
) (*GetAllProductReviewsResult, error) {

	offset := (input.Page - 1) * input.Limit

	// Step 1: Validate product
	product, err := s.Deps.Repo.GetProductByID(ctx, input.ProductID)
	if err != nil {
		return nil, errors.NewNotFoundError("product")
	}
	if product.IsBanned || product.IsArchived || !product.IsApproved {
		return nil, errors.NewValidationError("product", "not approved, banned, or archived")
	}

	// Step 2: Get all valid reviews for this product
	reviews, err := s.Deps.Repo.GetAllReviewsByProductID(ctx, sqlc.GetAllReviewsByProductIDParams{
		ProductID: input.ProductID,
		Limit:     int32(input.Limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, errors.NewTableError("product_reviews.select", err.Error())
	}

	// Step 3: Extract review IDs for batch reply fetch
	reviewIDs := make([]uuid.UUID, 0, len(reviews))
	for _, r := range reviews {
		reviewIDs = append(reviewIDs, r.ID)
	}

	// repliesMap := make(map[uuid.UUID]sqlc.ProductReviewReply)
	repliesMap := make(map[uuid.UUID]sqlc.GetRepliesByReviewIDsRow)

	if len(reviewIDs) > 0 {
		replies, err := s.Deps.Repo.GetRepliesByReviewIDs(ctx, reviewIDs)
		if err != nil {
			return nil, errors.NewTableError("product_review_replies.select", err.Error())
		}
		for _, reply := range replies {
			// repliesMap := make(map[uuid.UUID]sqlc.GetRepliesByReviewIDsRow)

			repliesMap[reply.ReviewID] = reply
		}
	}

	// Step 4: Format final response
	var formatted []ReviewWithReply
	for _, r := range reviews {
		review := ReviewWithReply{
			ReviewID:       r.ID,
			ReviewerUserID: r.ReviewerUserID,
			ReviewText:     r.ReviewText,
			// ReviewImageURL: r.ReviewImageUrl.Ptr(),
			ReviewImageURL: sqlnull.ToStringPtr(r.ReviewImageUrl),

			IsEdited:  r.IsEdited,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}

		if reply, ok := repliesMap[r.ID]; ok {
			review.Reply = &ReplyResponse{
				ReplyID:      reply.ID,
				SellerUserID: reply.SellerUserID,
				ReplyText:    reply.ReplyText,
				// ReplyImageURL: reply.ReplyImageUrl.Ptr(),
				ReplyImageURL: sqlnull.ToStringPtr(reply.ReplyImageUrl),

				IsEdited:  reply.IsEdited,
				CreatedAt: reply.CreatedAt,
				UpdatedAt: reply.UpdatedAt,
			}
		}

		formatted = append(formatted, review)
	}

	return &GetAllProductReviewsResult{
		Page:       input.Page,
		PerPage:    input.Limit,
		TotalItems: int(len(reviews)),
		Items:      formatted,
	}, nil
}
