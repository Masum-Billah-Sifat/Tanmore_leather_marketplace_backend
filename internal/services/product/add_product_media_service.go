// ------------------------------------------------------------
// 📁 File: internal/services/product/add_product_media_service.go
// 🧠 Handles adding media (image or promo video) to a product.
//     - Validates seller
//     - Validates product ownership
//     - If promo video, checks for duplicate
//     - Inserts media row into product_medias table
//     - Emits product.image.added or product.promo_video.added event

package product

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_add_media"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type AddProductMediaInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	MediaURL  string
	MediaType string // "image" or "promo_video"
}

// ------------------------------------------------------------
// 📤 Result to return
type AddProductMediaResult struct {
	ProductID uuid.UUID
	MediaID   uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type AddProductMediaServiceDeps struct {
	Repo repo.ProductAddMediaRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type AddProductMediaService struct {
	Deps AddProductMediaServiceDeps
}

// 🚀 Constructor
func NewAddProductMediaService(deps AddProductMediaServiceDeps) *AddProductMediaService {
	return &AddProductMediaService{Deps: deps}
}

// 🚀 Entrypoint
func (s *AddProductMediaService) Start(
	ctx context.Context,
	input AddProductMediaInput,
) (*AddProductMediaResult, error) {

	now := timeutil.NowUTC()
	var mediaID uuid.UUID

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate seller
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
		if !user.IsSellerProfileCreated || !user.IsSellerProfileApproved {
			return errors.NewValidationError("seller", "profile not approved or not created")
		}

		// ------------------------------------------------------------
		// Step 2: Validate product ownership
		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsArchived || product.IsBanned {
			return errors.NewValidationError("product", "banned or archived product")
		}

		// ------------------------------------------------------------
		// Step 3: Check if promo video already exists
		if input.MediaType == "video" {
			_, err := q.GetPromoVideoByProductID(ctx, sqlc.GetPromoVideoByProductIDParams{
				ProductID:  input.ProductID,
				MediaType:  "video",
				IsArchived: false,
			})
			if err == nil {
				return errors.NewConflictError("Promo video already exists for this product")
			}
			// If err != nil, that’s okay (no existing video)
		}

		// ------------------------------------------------------------
		// Step 4: Insert media row
		mediaUUID := uuidutil.New()
		mediaID = mediaUUID

		_, err = q.InsertProductMedia(ctx, sqlc.InsertProductMediaParams{
			ID:         mediaUUID,
			ProductID:  input.ProductID,
			MediaType:  input.MediaType,
			MediaUrl:   input.MediaURL,
			IsPrimary:  false,
			IsArchived: false,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		if err != nil {
			return errors.NewTableError("product_medias.insert", err.Error())
		}

		// ------------------------------------------------------------
		// Step 5: Emit event
		eventType := ""
		if input.MediaType == "image" {
			eventType = "product.image.added"
		} else {
			eventType = "product.promo_video.added"
		}

		payload := map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": input.ProductID,
			"media_url":  input.MediaURL,
			"media_id":   mediaUUID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    eventType,
			EventPayload: payloadBytes,
			DispatchedAt: sqlnull.TimePtr(nil),
			CreatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("events.insert", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AddProductMediaResult{
		ProductID: input.ProductID,
		MediaID:   mediaID,
		Status:    "media_added",
	}, nil
}
