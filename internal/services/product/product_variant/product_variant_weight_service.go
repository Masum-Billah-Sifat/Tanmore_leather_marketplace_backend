// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/update_variant_weight_service.go
// 🧠 Handles updating the weight (grams) of a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Updates the weight_grams
//     - Emits variant.weight.updated event
//     - Returns variant_id and updated weight

package product_variant

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_weight"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateVariantWeightInput struct {
	UserID      uuid.UUID
	ProductID   uuid.UUID
	VariantID   uuid.UUID
	WeightGrams int64
}

// 📤 Result to return
type UpdateVariantWeightResult struct {
	VariantID   uuid.UUID
	WeightGrams int64
	Status      string
}

// 🧱 Dependencies
type UpdateVariantWeightServiceDeps struct {
	Repo repo.ProductVariantUpdateWeightRepoInterface
}

// 🛠️ Service Definition
type UpdateVariantWeightService struct {
	Deps UpdateVariantWeightServiceDeps
}

// 🚀 Constructor
func NewUpdateVariantWeightService(deps UpdateVariantWeightServiceDeps) *UpdateVariantWeightService {
	return &UpdateVariantWeightService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateVariantWeightService) Start(
	ctx context.Context,
	input UpdateVariantWeightInput,
) (*UpdateVariantWeightResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Fetch snapshot (seller + product + variant + category)
		snapshot, err := q.GetVariantSnapshot(ctx, sqlc.GetVariantSnapshotParams{
			Sellerid:  input.UserID,
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		// ------------------------------------------------------------
		// Step 2: Check moderation rules
		if snapshot.Issellerarchived || snapshot.Issellerbanned {
			return errors.NewAuthError("seller is not allowed to modify products")
		}
		if snapshot.Isproductarchived || snapshot.Isproductbanned {
			return errors.NewValidationError("product", "cannot modify banned or archived product")
		}
		if snapshot.Iscategoryarchived {
			return errors.NewValidationError("category", "product's category is archived")
		}
		if snapshot.Isvariantarchived {
			return errors.NewValidationError("variant", "variant is archived")
		}

		// ------------------------------------------------------------
		if input.WeightGrams > int64(^int32(0)) || input.WeightGrams < int64(-1<<31) {
			return errors.NewValidationError("weight_grams", "value out of int32 range")
		}

		// Step 3: Update weight_grams
		err = q.UpdateVariantWeight(ctx, sqlc.UpdateVariantWeightParams{
			WeightGrams: int32(input.WeightGrams),
			UpdatedAt:   now,
			ID:          input.VariantID,
			ProductID:   input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.weight.updated event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":      input.UserID,
			"product_id":   input.ProductID,
			"variant_id":   input.VariantID,
			"weight_grams": input.WeightGrams,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.weight.updated",
			EventPayload: payload,
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

	return &UpdateVariantWeightResult{
		VariantID:   input.VariantID,
		WeightGrams: input.WeightGrams,
		Status:      "weight_updated",
	}, nil
}
