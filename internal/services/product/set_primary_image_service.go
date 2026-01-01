// ------------------------------------------------------------
// 📁 File: internal/services/product/set_primary_image_service.go
// 🧠 Handles setting an existing image as the primary image.
//     - Validates seller
//     - Validates product ownership
//     - Unsets existing primary images
//     - Sets selected image as primary
//     - Emits product.image.set_primary event

package product

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_set_primary_image"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type SetPrimaryImageInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	MediaID   uuid.UUID
}

// ------------------------------------------------------------
// 📤 Result to return
type SetPrimaryImageResult struct {
	ProductID      uuid.UUID
	PrimaryImageID uuid.UUID
	Status         string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type SetPrimaryImageServiceDeps struct {
	Repo repo.ProductSetPrimaryImageRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type SetPrimaryImageService struct {
	Deps SetPrimaryImageServiceDeps
}

// 🚀 Constructor
func NewSetPrimaryImageService(deps SetPrimaryImageServiceDeps) *SetPrimaryImageService {
	return &SetPrimaryImageService{Deps: deps}
}

// 🚀 Entrypoint
func (s *SetPrimaryImageService) Start(
	ctx context.Context,
	input SetPrimaryImageInput,
) (*SetPrimaryImageResult, error) {
	now := timeutil.NowUTC()

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
		// Step 3: Fetch media row
		media, err := q.GetProductMediaByID(ctx, sqlc.GetProductMediaByIDParams{
			ID:        input.MediaID,
			ProductID: input.ProductID,
			MediaType: "image", // ensure type match
		})
		if err != nil {
			return errors.NewNotFoundError("media")
		}
		if media.IsArchived {
			return errors.NewValidationError("media", "cannot set archived image as primary")
		}
		if media.IsPrimary {
			return errors.NewConflictError("media")
		}

		// ------------------------------------------------------------
		// Step 3: Unset all previous is_primary flags
		err = q.UnsetAllPrimaryImages(ctx, sqlc.UnsetAllPrimaryImagesParams{
			ProductID:  input.ProductID,
			MediaType:  "image",
			IsArchived: false,
			IsPrimary:  false, // New value
		})
		if err != nil {
			return errors.NewTableError("product_medias.unset_primary", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Set selected image as primary
		err = q.SetAsPrimaryImage(ctx, sqlc.SetAsPrimaryImageParams{
			ID:        input.MediaID,
			ProductID: input.ProductID,
			MediaType: "image",
			IsPrimary: true,
		})
		if err != nil {
			return errors.NewTableError("product_medias.set_primary", err.Error())
		}

		// ------------------------------------------------------------
		// Step 5: Emit event
		payload := map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": input.ProductID,
			"media_id":   input.MediaID,
			"media_url":  media.MediaUrl, // ✅ include it

		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "product.image.set_primary",
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

	return &SetPrimaryImageResult{
		ProductID:      input.ProductID,
		PrimaryImageID: input.MediaID,
		Status:         "primary_image_set",
	}, nil
}
