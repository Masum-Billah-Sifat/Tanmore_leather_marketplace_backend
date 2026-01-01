// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/remove_variant_retail_discount_service.go
// 🧠 Handles removing a retail discount from a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Checks if retail discount exists
//     - Removes the discount fields from DB
//     - Emits variant.retail_discount.removed event
//     - Returns variant_id and status

package product_variant

import (
	"context"
	"encoding/json"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_remove_retail_discount"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input from handler

type RemoveVariantRetailDiscountInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	VariantID uuid.UUID
}

// ------------------------------------------------------------
// 👋 Result to return

type RemoveVariantRetailDiscountResult struct {
	VariantID uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🔧 Dependencies

type RemoveVariantRetailDiscountServiceDeps struct {
	Repo repo.ProductVariantRemoveDiscountRepoInterface
}

// 🔧 Service Definition
type RemoveVariantRetailDiscountService struct {
	Deps RemoveVariantRetailDiscountServiceDeps
}

// 🚀 Constructor
func NewRemoveVariantRetailDiscountService(deps RemoveVariantRetailDiscountServiceDeps) *RemoveVariantRetailDiscountService {
	return &RemoveVariantRetailDiscountService{Deps: deps}
}

// 🚀 Entrypoint
func (s *RemoveVariantRetailDiscountService) Start(
	ctx context.Context,
	input RemoveVariantRetailDiscountInput,
) (*RemoveVariantRetailDiscountResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Fetch snapshot
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
		if !snapshot.Hasretaildiscount {
			return errors.NewConflictError("variant has no retail discount to remove")
		}

		// ------------------------------------------------------------
		// Step 3: Remove retail discount fields
		err = q.DisableRetailDiscount(ctx, sqlc.DisableRetailDiscountParams{
			Hasretaildiscount:  false,
			Retaildiscounttype: sqlnull.String(""), // set to NULL
			Retaildiscount:     sqlnull.Int64(0),   // set to NULL
			UpdatedAt:          now,
			ID:                 input.VariantID,
			ProductID:          input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.retail_discount.removed event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": input.ProductID,
			"variant_id": input.VariantID,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.retail_discount.removed",
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

	return &RemoveVariantRetailDiscountResult{
		VariantID: input.VariantID,
		Status:    "retail_discount_removed",
	}, nil
}
