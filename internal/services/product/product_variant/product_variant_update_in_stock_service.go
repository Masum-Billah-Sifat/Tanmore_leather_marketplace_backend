// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/update_variant_in_stock_service.go
// 🧠 Handles updating the in_stock status of a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Skips if status is unchanged
//     - Updates the in_stock flag
//     - Emits variant.in_stock.updated event
//     - Returns variant_id and new in_stock status

package product_variant

import (
	"context"
	"encoding/json"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_in_stock"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateVariantInStockInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	VariantID uuid.UUID
	InStock   bool
}

// ------------------------------------------------------------
// 📤 Result to return
type UpdateVariantInStockResult struct {
	VariantID uuid.UUID
	InStock   bool
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateVariantInStockServiceDeps struct {
	Repo repo.ProductVariantUpdateInStockRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type UpdateVariantInStockService struct {
	Deps UpdateVariantInStockServiceDeps
}

// 🚀 Constructor
func NewUpdateVariantInStockService(deps UpdateVariantInStockServiceDeps) *UpdateVariantInStockService {
	return &UpdateVariantInStockService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateVariantInStockService) Start(
	ctx context.Context,
	input UpdateVariantInStockInput,
) (*UpdateVariantInStockResult, error) {

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
		// Step 3: Skip update if same value
		if snapshot.Isvariantinstock == input.InStock {
			return errors.NewValidationError("in_stock", "value is already set to the requested state")
		}

		// ------------------------------------------------------------
		// Step 4: Update in_stock field
		err = q.UpdateVariantInStock(ctx, sqlc.UpdateVariantInStockParams{
			InStock:   input.InStock,
			UpdatedAt: now,
			ID:        input.VariantID,
			ProductID: input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 5: Emit variant.in_stock.updated event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": input.ProductID,
			"variant_id": input.VariantID,
			"in_stock":   input.InStock,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.in_stock.updated",
			EventPayload: payload,
			DispatchedAt: sqlnull.Time(time.Time{}),
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

	return &UpdateVariantInStockResult{
		VariantID: input.VariantID,
		InStock:   input.InStock,
		Status:    "in_stock_status_updated",
	}, nil
}
