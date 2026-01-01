// ------------------------------------------------------------
// 📁 File: internal/services/product/product_variant/update_variant_retail_discount_service.go
// 🧠 Handles updating the retail discount of a variant.
//     - Validates seller, product, and variant using snapshot
//     - Confirms that retail discount is enabled
//     - Performs COALESCE-based update for partial fields
//     - Emits variant.retail_discount.updated event
//     - Returns updated fields and variant ID

package product_variant

import (
	"context"
	"database/sql"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_update_discount"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	// "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateVariantRetailDiscountInput struct {
	UserID             uuid.UUID
	ProductID          uuid.UUID
	VariantID          uuid.UUID
	RetailDiscount     *int64  // optional
	RetailDiscountType *string // optional
}

// ------------------------------------------------------------
// 📤 Output to handler
type UpdateVariantRetailDiscountResult struct {
	VariantID     uuid.UUID
	UpdatedFields map[string]interface{}
	Status        string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateVariantRetailDiscountServiceDeps struct {
	Repo repo.ProductVariantUpdateDiscountRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service definition
type UpdateVariantRetailDiscountService struct {
	Deps UpdateVariantRetailDiscountServiceDeps
}

// 🚀 Constructor
func NewUpdateVariantRetailDiscountService(deps UpdateVariantRetailDiscountServiceDeps) *UpdateVariantRetailDiscountService {
	return &UpdateVariantRetailDiscountService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateVariantRetailDiscountService) Start(
	ctx context.Context,
	input UpdateVariantRetailDiscountInput,
) (*UpdateVariantRetailDiscountResult, error) {

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

		if snapshot.Issellerarchived || snapshot.Issellerbanned || !snapshot.Issellerapproved {
			return errors.NewAuthError("seller is not allowed to update")
		}
		if snapshot.Isproductarchived || snapshot.Isproductbanned || !snapshot.Isproductapproved {
			return errors.NewValidationError("product", "product is not in updatable state")
		}
		if snapshot.Iscategoryarchived {
			return errors.NewValidationError("category", "cannot update variant under archived category")
		}
		if snapshot.Isvariantarchived {
			return errors.NewValidationError("variant", "variant is archived")
		}
		if !snapshot.Hasretaildiscount {
			return errors.NewValidationError("retail_discount", "retail discount is not enabled for this variant")
		}

		// ------------------------------------------------------------
		// Step 2: Update retail discount (COALESCE-style)

		var discountNull sql.NullInt64
		var discountTypeNull sql.NullString

		updatedFields := make(map[string]interface{})

		if input.RetailDiscount != nil {
			discountNull = sqlnull.Int64(*input.RetailDiscount)
			updatedFields["retail_discount"] = *input.RetailDiscount
		}

		if input.RetailDiscountType != nil {
			discountTypeNull = sqlnull.String(*input.RetailDiscountType)
			updatedFields["retail_discount_type"] = *input.RetailDiscountType
		}

		err = q.UpdateRetailDiscount(ctx, sqlc.UpdateRetailDiscountParams{
			Retaildiscount:     discountNull,     // ✅ sql.NullInt64 → retail_discount
			Retaildiscounttype: discountTypeNull, // ✅ sql.NullString → retail_discount_type
			UpdatedAt:          now,
			VariantID:          input.VariantID,
			ProductID:          input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

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
			EventType:    "variant.retail_discount.updated",
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

	// 🔁 Build response
	return &UpdateVariantRetailDiscountResult{
		VariantID:     input.VariantID,
		UpdatedFields: buildRetailDiscountFieldMap(input.RetailDiscount, input.RetailDiscountType),
		Status:        "retail_discount_updated",
	}, nil
}

// ------------------------------------------------------------
// 🔧 Build updated fields map
func buildRetailDiscountFieldMap(discount *int64, discountType *string) map[string]interface{} {
	m := make(map[string]interface{})
	if discount != nil {
		m["retail_discount"] = *discount
	}
	if discountType != nil {
		m["retail_discount_type"] = *discountType
	}
	return m
}
