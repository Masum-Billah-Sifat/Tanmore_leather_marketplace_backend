// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/update_variant_info_service.go
// 🧠 Handles updating color and/or size of a product variant.
//     - Validates seller & ownership via snapshot
//     - Updates non-null fields (COALESCE strategy)
//     - Emits variant.info.updated event
//     - Returns updated fields and variant ID

package product_variant

import (
	"context"
	"database/sql"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_info"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateVariantInfoInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	VariantID uuid.UUID
	Color     *string // Optional
	Size      *string // Optional
}

// ------------------------------------------------------------
// 📤 Result to return
type UpdateVariantInfoResult struct {
	VariantID     uuid.UUID
	UpdatedFields map[string]string
	Status        string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateVariantInfoServiceDeps struct {
	Repo repo.ProductVariantUpdateInfoRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type UpdateVariantInfoService struct {
	Deps UpdateVariantInfoServiceDeps
}

// 🚀 Constructor
func NewUpdateVariantInfoService(deps UpdateVariantInfoServiceDeps) *UpdateVariantInfoService {
	return &UpdateVariantInfoService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateVariantInfoService) Start(
	ctx context.Context,
	input UpdateVariantInfoInput,
) (*UpdateVariantInfoResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate snapshot (ownership + moderation)
		snapshot, err := q.GetVariantSnapshot(ctx, sqlc.GetVariantSnapshotParams{
			Sellerid:  input.UserID,
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		if !snapshot.Issellerapproved {
			return errors.NewValidationError("seller", "seller not approved")
		}
		if snapshot.Issellerarchived || snapshot.Issellerbanned {
			return errors.NewAuthError("seller is not allowed to modify products")
		}
		if snapshot.Isproductarchived || snapshot.Isproductbanned || !snapshot.Isproductapproved {
			return errors.NewValidationError("product", "cannot update non-approved, archived or banned product")
		}
		if snapshot.Isvariantarchived {
			return errors.NewValidationError("variant", "variant is archived")
		}

		// ------------------------------------------------------------
		// Step 2: Update color and/or size
		var colorNull sql.NullString
		var sizeNull sql.NullString

		updatedFields := make(map[string]string)

		if input.Color != nil {
			colorNull = sqlnull.String(*input.Color)
			updatedFields["color"] = *input.Color
		}

		if input.Size != nil {
			sizeNull = sqlnull.String(*input.Size)
			updatedFields["size"] = *input.Size
		}

		err = q.UpdateVariantColorSize(ctx, sqlc.UpdateVariantColorSizeParams{
			Color:     colorNull,
			Size:      sizeNull,
			UpdatedAt: now,
			VariantID: input.VariantID,
			ProductID: input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 3: Emit variant.info.updated event
		eventPayload := map[string]interface{}{
			"user_id":        input.UserID,
			"product_id":     input.ProductID,
			"variant_id":     input.VariantID,
			"updated_fields": updatedFields,
		}

		payloadBytes, err := json.Marshal(eventPayload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.info.updated",
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

	return &UpdateVariantInfoResult{
		VariantID:     input.VariantID,
		UpdatedFields: buildFieldMap(input.Color, input.Size),
		Status:        "variant_info_updated",
	}, nil
}

// 🔧 Helper to return only updated fields
func buildFieldMap(color *string, size *string) map[string]string {
	m := make(map[string]string)
	if color != nil {
		m["color"] = *color
	}
	if size != nil {
		m["size"] = *size
	}
	return m
}
