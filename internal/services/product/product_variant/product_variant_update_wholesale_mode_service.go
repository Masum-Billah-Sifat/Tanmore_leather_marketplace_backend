// ------------------------------------------------------------
// 📁 File: internal/services/product/product_variant/edit_variant_wholesale_info_service.go
// 🧠 Handles editing the wholesale info of a variant.
//     - Validates seller, product, and variant using snapshot
//     - Confirms that wholesale mode is enabled
//     - Performs COALESCE-based update for partial fields
//     - Emits variant.wholesale_mode.updated event
//     - Returns updated fields and variant ID

package product_variant

import (
	"context"
	"database/sql"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_wholesale_mode"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type EditWholesaleInfoInput struct {
	UserID          uuid.UUID
	ProductID       uuid.UUID
	VariantID       uuid.UUID
	WholesalePrice  *int64 // optional
	MinQtyWholesale *int32 // optional
}

// ------------------------------------------------------------
// 📤 Output to handler
type EditWholesaleInfoResult struct {
	VariantID     uuid.UUID
	UpdatedFields map[string]interface{}
	Status        string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type EditWholesaleInfoServiceDeps struct {
	Repo repo.ProductVariantEditWholesaleInfoRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service definition
type EditWholesaleInfoService struct {
	Deps EditWholesaleInfoServiceDeps
}

// 🚀 Constructor
func NewEditWholesaleInfoService(deps EditWholesaleInfoServiceDeps) *EditWholesaleInfoService {
	return &EditWholesaleInfoService{Deps: deps}
}

// 🚀 Entrypoint
func (s *EditWholesaleInfoService) Start(
	ctx context.Context,
	input EditWholesaleInfoInput,
) (*EditWholesaleInfoResult, error) {
	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate snapshot
		snapshot, err := q.GetVariantSnapshot(ctx, sqlc.GetVariantSnapshotParams{
			Sellerid:  input.UserID,
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		if snapshot.Issellerarchived || snapshot.Issellerbanned {
			return errors.NewAuthError("seller is not allowed to update")
		}
		if snapshot.Isproductarchived || snapshot.Isproductbanned {
			return errors.NewValidationError("product", "product is not in updatable state")
		}
		if snapshot.Iscategoryarchived {
			return errors.NewValidationError("category", "cannot update variant under archived category")
		}
		if snapshot.Isvariantarchived {
			return errors.NewValidationError("variant", "variant is archived")
		}
		if !snapshot.Haswholesaleenabled {
			return errors.NewValidationError("wholesale_mode", "wholesale mode is not enabled for this variant")
		}

		// ------------------------------------------------------------
		// Step 2: Update wholesale info (COALESCE-style)

		var priceNull sql.NullInt64
		var qtyNull sql.NullInt32

		updatedFields := make(map[string]interface{})

		if input.WholesalePrice != nil {
			priceNull = sqlnull.Int64(*input.WholesalePrice)
			updatedFields["wholesale_price"] = *input.WholesalePrice
		}
		if input.MinQtyWholesale != nil {
			qtyNull = sqlnull.Int32From32(*input.MinQtyWholesale)
			updatedFields["min_qty_wholesale"] = *input.MinQtyWholesale
		}

		err = q.UpdateWholesaleMode(ctx, sqlc.UpdateWholesaleModeParams{
			WholesalePrice:  priceNull,
			MinQtyWholesale: qtyNull,
			UpdatedAt:       now,
			VariantID:       input.VariantID,
			ProductID:       input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 3: Insert event
		payload := map[string]interface{}{
			"user_id":        input.UserID,
			"product_id":     input.ProductID,
			"variant_id":     input.VariantID,
			"updated_fields": updatedFields,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuid.New(),
			Userid:       input.UserID,
			EventType:    "variant.wholesale_mode.updated",
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

	// ------------------------------------------------------------
	// Final Output
	return &EditWholesaleInfoResult{
		VariantID:     input.VariantID,
		UpdatedFields: buildWholesaleUpdateFieldMap(input.WholesalePrice, input.MinQtyWholesale),
		Status:        "wholesale_mode_updated",
	}, nil
}

// ------------------------------------------------------------
// 🔧 Build updated fields map
func buildWholesaleUpdateFieldMap(price *int64, qty *int32) map[string]interface{} {
	m := make(map[string]interface{})
	if price != nil {
		m["wholesale_price"] = *price
	}
	if qty != nil {
		m["min_qty_wholesale"] = *qty
	}
	return m
}
