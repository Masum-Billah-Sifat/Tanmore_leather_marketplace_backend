// ------------------------------------------------------------
// 📁 File: internal/services/product/product_variant/update_variant_stock_quantity_service.go
// 🧠 Handles updating the stock quantity of a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Updates the stock_quantity
//     - Emits variant.stock_quantity.updated event
//     - Returns variant_id and updated quantity

package product_variant

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_stock_quantity"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateVariantStockQuantityInput struct {
	UserID        uuid.UUID
	ProductID     uuid.UUID
	VariantID     uuid.UUID
	StockQuantity int64
}

// ------------------------------------------------------------
// 📤 Result to return
type UpdateVariantStockQuantityResult struct {
	VariantID     uuid.UUID
	StockQuantity int64
	Status        string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateVariantStockQuantityServiceDeps struct {
	Repo repo.ProductVariantUpdateStockQuantityRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type UpdateVariantStockQuantityService struct {
	Deps UpdateVariantStockQuantityServiceDeps
}

// 🚀 Constructor
func NewUpdateVariantStockQuantityService(deps UpdateVariantStockQuantityServiceDeps) *UpdateVariantStockQuantityService {
	return &UpdateVariantStockQuantityService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateVariantStockQuantityService) Start(
	ctx context.Context,
	input UpdateVariantStockQuantityInput,
) (*UpdateVariantStockQuantityResult, error) {

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
		if input.StockQuantity > int64(^int32(0)) {
			return errors.NewValidationError("stock_quantity", "value too large")
		}

		// Step 3: Update stock quantity
		err = q.UpdateVariantStockQuantity(ctx, sqlc.UpdateVariantStockQuantityParams{
			StockQuantity: int32(input.StockQuantity),
			UpdatedAt:     now,
			ID:            input.VariantID,
			ProductID:     input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.stock_quantity.updated event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":        input.UserID,
			"product_id":     input.ProductID,
			"variant_id":     input.VariantID,
			"stock_quantity": input.StockQuantity,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.stock_quantity.updated",
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

	return &UpdateVariantStockQuantityResult{
		VariantID:     input.VariantID,
		StockQuantity: input.StockQuantity,
		Status:        "stock_quantity_updated",
	}, nil
}
