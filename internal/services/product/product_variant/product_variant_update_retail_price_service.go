// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/update_variant_retail_price_service.go
// 🧠 Handles updating the retail price of a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Updates the retail_price
//     - Emits variant.retail_price.updated event
//     - Returns variant_id and updated price

package product_variant

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_price"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateVariantRetailPriceInput struct {
	UserID      uuid.UUID
	ProductID   uuid.UUID
	VariantID   uuid.UUID
	RetailPrice int64
}

// ------------------------------------------------------------

// 📤 Result to return
type UpdateVariantRetailPriceResult struct {
	VariantID   uuid.UUID
	RetailPrice int64
	Status      string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateVariantRetailPriceServiceDeps struct {
	Repo repo.ProductVariantUpdatePriceRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type UpdateVariantRetailPriceService struct {
	Deps UpdateVariantRetailPriceServiceDeps
}

// 🚀 Constructor
func NewUpdateVariantRetailPriceService(deps UpdateVariantRetailPriceServiceDeps) *UpdateVariantRetailPriceService {
	return &UpdateVariantRetailPriceService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateVariantRetailPriceService) Start(
	ctx context.Context,
	input UpdateVariantRetailPriceInput,
) (*UpdateVariantRetailPriceResult, error) {

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
		// Step 3: Update retail price
		err = q.UpdateVariantRetailPrice(ctx, sqlc.UpdateVariantRetailPriceParams{
			RetailPrice: input.RetailPrice,
			UpdatedAt:   now,
			ID:          input.VariantID,
			ProductID:   input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.retail_price.updated event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":      input.UserID,
			"product_id":   input.ProductID,
			"variant_id":   input.VariantID,
			"retail_price": input.RetailPrice,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.retail_price.updated",
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

	return &UpdateVariantRetailPriceResult{
		VariantID:   input.VariantID,
		RetailPrice: input.RetailPrice,
		Status:      "retail_price_updated",
	}, nil
}
