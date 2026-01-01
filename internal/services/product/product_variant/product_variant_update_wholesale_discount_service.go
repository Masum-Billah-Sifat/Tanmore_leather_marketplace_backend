// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/update_variant_wholesale_discount_service.go
// 🧠 Handles updating the wholesale discount of a variant.
//     - Validates seller, product, variant, and category via snapshot
//     - Confirms wholesale mode is already enabled
//     - Updates discount fields using COALESCE
//     - Emits variant.wholesale_discount.updated event
//     - Returns updated fields and variant ID

package product_variant

import (
	"context"
	"database/sql"
	"encoding/json"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_wholesale_discount"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateWholesaleDiscountInput struct {
	UserID                uuid.UUID
	ProductID             uuid.UUID
	VariantID             uuid.UUID
	WholesaleDiscount     *int64  // optional
	WholesaleDiscountType *string // optional
}

// ------------------------------------------------------------
// 📤 Output to handler
type UpdateWholesaleDiscountResult struct {
	VariantID     uuid.UUID
	UpdatedFields map[string]interface{}
	Status        string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateWholesaleDiscountServiceDeps struct {
	Repo repo.ProductVariantUpdateWholesaleDiscountRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service definition
type UpdateWholesaleDiscountService struct {
	Deps UpdateWholesaleDiscountServiceDeps
}

// 🚀 Constructor
func NewUpdateWholesaleDiscountService(deps UpdateWholesaleDiscountServiceDeps) *UpdateWholesaleDiscountService {
	return &UpdateWholesaleDiscountService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateWholesaleDiscountService) Start(
	ctx context.Context,
	input UpdateWholesaleDiscountInput,
) (*UpdateWholesaleDiscountResult, error) {
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
		// Step 2: Build COALESCE update payload
		var discountNull sql.NullInt64
		var typeNull sql.NullString

		updatedFields := make(map[string]interface{})

		if input.WholesaleDiscount != nil {
			discountNull = sqlnull.Int64(*input.WholesaleDiscount)
			updatedFields["wholesale_discount"] = *input.WholesaleDiscount
		}
		if input.WholesaleDiscountType != nil {
			typeNull = sqlnull.String(*input.WholesaleDiscountType)
			updatedFields["wholesale_discount_type"] = *input.WholesaleDiscountType
		}

		if len(updatedFields) == 0 {
			return errors.NewValidationError("payload", "at least one field is required")
		}

		err = q.UpdateWholesaleDiscount(ctx, sqlc.UpdateWholesaleDiscountParams{
			Wholesalediscount:     discountNull,
			Wholesalediscounttype: typeNull,
			UpdatedAt:             now,
			VariantID:             input.VariantID,
			ProductID:             input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 3: Emit event
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
			EventType:    "variant.wholesale_discount.updated",
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

	return &UpdateWholesaleDiscountResult{
		VariantID:     input.VariantID,
		UpdatedFields: buildUpdatedDiscountFieldMap(input.WholesaleDiscount, input.WholesaleDiscountType),
		Status:        "wholesale_discount_updated",
	}, nil
}

// ------------------------------------------------------------
// 🔧 Build updated fields map
func buildUpdatedDiscountFieldMap(discount *int64, discountType *string) map[string]interface{} {
	m := make(map[string]interface{})
	if discount != nil {
		m["wholesale_discount"] = *discount
	}
	if discountType != nil {
		m["wholesale_discount_type"] = *discountType
	}
	return m
}
