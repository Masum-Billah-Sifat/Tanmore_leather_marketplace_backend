// ------------------------------------------------------------
// 📁 File: internal/services/product/archive_product_media_service.go
// 🧠 Handles archiving media (image or promo video) from a product.
//     - Validates seller
//     - Validates product ownership
//     - Validates media existence
//     - Prevents archiving if last image
//     - Archives media row
//     - Emits product.image.removed or product.promo_video.removed event

package product

import (
	"context"
	"encoding/json"
	"log"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_archive_media"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type ArchiveProductMediaInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	MediaID   uuid.UUID
	MediaType string // "image" or "promo_video"
}

// ------------------------------------------------------------
// 📤 Result to return
type ArchiveProductMediaResult struct {
	ProductID uuid.UUID
	MediaID   uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type ArchiveProductMediaServiceDeps struct {
	Repo repo.ProductArchiveMediaRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type ArchiveProductMediaService struct {
	Deps ArchiveProductMediaServiceDeps
}

// 🚀 Constructor
func NewArchiveProductMediaService(deps ArchiveProductMediaServiceDeps) *ArchiveProductMediaService {
	return &ArchiveProductMediaService{Deps: deps}
}

// // 🚀 Entrypoint
// func (s *ArchiveProductMediaService) Start(
// 	ctx context.Context,
// 	input ArchiveProductMediaInput,
// ) (*ArchiveProductMediaResult, error) {

// 	now := timeutil.NowUTC()

// 	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// 		// ------------------------------------------------------------
// 		// Step 1: Validate seller
// 		user, err := q.GetUserByID(ctx, input.UserID)
// 		if err != nil {
// 			return errors.NewNotFoundError("user")
// 		}
// 		if user.IsArchived {
// 			return errors.NewAuthError("user is archived")
// 		}
// 		if user.IsBanned {
// 			return errors.NewAuthError("user is banned")
// 		}
// 		if !user.IsSellerProfileCreated || !user.IsSellerProfileApproved {
// 			return errors.NewValidationError("seller", "profile not approved or not created")
// 		}

// 		// ------------------------------------------------------------
// 		// Step 2: Validate product ownership
// 		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
// 			ID:       input.ProductID,
// 			SellerID: input.UserID,
// 		})
// 		if err != nil {
// 			return errors.NewNotFoundError("product")
// 		}
// 		if product.IsArchived || product.IsBanned {
// 			return errors.NewValidationError("product", "banned or archived product")
// 		}

// 		// ------------------------------------------------------------
// 		// Step 3: Fetch media row
// 		media, err := q.GetProductMediaByID(ctx, sqlc.GetProductMediaByIDParams{
// 			ID:        input.MediaID,
// 			ProductID: input.ProductID,
// 			MediaType: input.MediaType,
// 		})
// 		if err != nil {
// 			return errors.NewNotFoundError("media")
// 		}

// 		if media.IsArchived {
// 			return errors.NewConflictError("Media already archived")
// 		}

// 		// here we have to check one more thing because if the media of type image and it's set as primary you cant delete it as well? clear
// 		// hold on we need slight update here

// 		// 🛑 Prevent archiving last active image
// 		if input.MediaType == "image" {
// 			count, err := q.CountActiveImagesForProduct(ctx, sqlc.CountActiveImagesForProductParams{
// 				ProductID:  input.ProductID,
// 				MediaType:  "image",
// 				IsArchived: false,
// 			})
// 			if err != nil {
// 				return errors.NewServerError("count images")
// 			}
// 			if count == 1 {
// 				return errors.NewValidationError("media", "cannot archive the last active image")
// 			}
// 		}

// 		// ------------------------------------------------------------
// 		// Step 4: Archive media
// 		err = q.ArchiveProductMedia(ctx, sqlc.ArchiveProductMediaParams{
// 			ID:         input.MediaID,
// 			ProductID:  input.ProductID,
// 			MediaType:  input.MediaType,
// 			IsArchived: true,
// 		})
// 		if err != nil {
// 			return errors.NewTableError("product_medias.archive", err.Error())
// 		}

// 		// ------------------------------------------------------------
// 		// Step 5: Emit event
// 		eventType := ""
// 		if input.MediaType == "image" {
// 			eventType = "product.image.removed"
// 		} else {
// 			eventType = "product.promo_video.removed"
// 		}

// 		payload := map[string]interface{}{
// 			"user_id":    input.UserID,
// 			"product_id": input.ProductID,
// 			"media_id":   input.MediaID,
// 			"media_url":  media.MediaUrl, // ✅ required

// 		}

// 		payloadBytes, err := json.Marshal(payload)
// 		if err != nil {
// 			return errors.NewServerError("marshal event payload")
// 		}

// 		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
// 			ID:           uuidutil.New(),
// 			Userid:       input.UserID,
// 			EventType:    eventType,
// 			EventPayload: payloadBytes,
// 			DispatchedAt: sqlnull.TimePtr(nil),
// 			CreatedAt:    now,
// 		})
// 		if err != nil {
// 			return errors.NewTableError("events.insert", err.Error())
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	status := "media_archived"
// 	switch input.MediaType {
// 	case "video":
// 		status = "promo_video_removed"
// 	case "image":
// 		status = "image_archived"
// 	}

// 	return &ArchiveProductMediaResult{
// 		ProductID: input.ProductID,
// 		MediaID:   input.MediaID,
// 		Status:    status,
// 	}, nil
// }

func (s *ArchiveProductMediaService) Start(
	ctx context.Context,
	input ArchiveProductMediaInput,
) (*ArchiveProductMediaResult, error) {

	log.Println("🟡 [ArchiveMedia] START",
		"user_id=", input.UserID,
		"product_id=", input.ProductID,
		"media_id=", input.MediaID,
		"media_type=", input.MediaType,
	)

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {

		// ------------------------------------------------------------
		log.Println("🟢 [ArchiveMedia] Step 1: Validate seller")

		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			log.Println("❌ [ArchiveMedia] User not found", err)
			return errors.NewNotFoundError("user")
		}

		if user.IsArchived {
			log.Println("❌ [ArchiveMedia] User archived")
			return errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			log.Println("❌ [ArchiveMedia] User banned")
			return errors.NewAuthError("user is banned")
		}
		if !user.IsSellerProfileCreated || !user.IsSellerProfileApproved {
			log.Println("❌ [ArchiveMedia] Seller profile not approved")
			return errors.NewValidationError("seller", "profile not approved or not created")
		}

		// ------------------------------------------------------------
		log.Println("🟢 [ArchiveMedia] Step 2: Validate product ownership")

		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.UserID,
		})
		if err != nil {
			log.Println("❌ [ArchiveMedia] Product not found", err)
			return errors.NewNotFoundError("product")
		}

		if product.IsArchived || product.IsBanned {
			log.Println("❌ [ArchiveMedia] Product archived or banned")
			return errors.NewValidationError("product", "banned or archived product")
		}

		// ------------------------------------------------------------
		log.Println("🟢 [ArchiveMedia] Step 3: Fetch media row")

		media, err := q.GetProductMediaByID(ctx, sqlc.GetProductMediaByIDParams{
			ID:        input.MediaID,
			ProductID: input.ProductID,
			MediaType: input.MediaType,
		})
		if err != nil {
			log.Println("❌ [ArchiveMedia] Media not found", err)
			return errors.NewNotFoundError("media")
		}

		log.Println("ℹ️ [ArchiveMedia] Media found",
			"is_archived=", media.IsArchived,
			"is_primary=", media.IsPrimary,
			"media_url=", media.MediaUrl,
		)

		if media.IsArchived {
			log.Println("❌ [ArchiveMedia] Media already archived")
			return errors.NewConflictError("Media already archived")
		}

		// ------------------------------------------------------------
		log.Println("🟢 [ArchiveMedia] Step 4: Prevent last active image deletion")

		if input.MediaType == "image" {
			count, err := q.CountActiveImagesForProduct(ctx, sqlc.CountActiveImagesForProductParams{
				ProductID:  input.ProductID,
				MediaType:  "image",
				IsArchived: false,
			})
			if err != nil {
				log.Println("❌ [ArchiveMedia] Failed counting images", err)
				return errors.NewServerError("count images")
			}

			log.Println("ℹ️ [ArchiveMedia] Active image count =", count)

			if count == 1 {
				log.Println("❌ [ArchiveMedia] Attempt to delete last active image")
				return errors.NewValidationError("media", "cannot archive the last active image")
			}
		}

		// ------------------------------------------------------------
		log.Println("🟢 [ArchiveMedia] Step 5: Archiving media")

		err = q.ArchiveProductMedia(ctx, sqlc.ArchiveProductMediaParams{
			ID:         input.MediaID,
			ProductID:  input.ProductID,
			MediaType:  input.MediaType,
			IsArchived: true,
		})
		if err != nil {
			log.Println("❌ [ArchiveMedia] Failed to archive media", err)
			return errors.NewTableError("product_medias.archive", err.Error())
		}

		// ------------------------------------------------------------
		log.Println("🟢 [ArchiveMedia] Step 6: Emit event")

		eventType := "product.image.removed"
		if input.MediaType == "video" {
			eventType = "product.promo_video.removed"
		}

		payload := map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": input.ProductID,
			"media_id":   input.MediaID,
			"media_url":  media.MediaUrl,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Println("❌ [ArchiveMedia] Failed to marshal event payload", err)
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
			log.Println("❌ [ArchiveMedia] Failed inserting event", err)
			return errors.NewTableError("events.insert", err.Error())
		}

		log.Println("✅ [ArchiveMedia] SUCCESS")
		return nil
	})

	if err != nil {
		log.Println("🔥 [ArchiveMedia] TX FAILED:", err)
		return nil, err
	}

	return &ArchiveProductMediaResult{
		ProductID: input.ProductID,
		MediaID:   input.MediaID,
		Status:    "media_archived",
	}, nil
}
