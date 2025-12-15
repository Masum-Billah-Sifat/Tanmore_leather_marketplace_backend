// ------------------------------------------------------------
// 📁 File: internal/services/product/product_variant/remove_variant_wholesale_discount_service.go
// 🧠 Handles removing a wholesale discount from a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Checks if wholesale discount exists
//     - Removes wholesale discount fields from DB
//     - Emits variant.wholesale_discount.removed event
//     - Returns variant_id and status

package product_variant

import (
	"context"
	"encoding/json"
	"time"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_disable_wholesale_discount"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input from handler

type RemoveVariantWholesaleDiscountInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	VariantID uuid.UUID
}

// ------------------------------------------------------------
// 👋 Result to return

type RemoveVariantWholesaleDiscountResult struct {
	VariantID uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🔧 Dependencies

type RemoveVariantWholesaleDiscountServiceDeps struct {
	Repo repo.ProductVariantRemoveWholesaleDiscountRepoInterface
}

// 🔧 Service Definition
type RemoveVariantWholesaleDiscountService struct {
	Deps RemoveVariantWholesaleDiscountServiceDeps
}

// 🚀 Constructor
func NewRemoveVariantWholesaleDiscountService(deps RemoveVariantWholesaleDiscountServiceDeps) *RemoveVariantWholesaleDiscountService {
	return &RemoveVariantWholesaleDiscountService{Deps: deps}
}

// 🚀 Entrypoint
func (s *RemoveVariantWholesaleDiscountService) Start(
	ctx context.Context,
	input RemoveVariantWholesaleDiscountInput,
) (*RemoveVariantWholesaleDiscountResult, error) {

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
		if !snapshot.Haswholesalediscount {
			return errors.NewConflictError("variant has no wholesale discount to remove")
		}

		// ------------------------------------------------------------
		// Step 3: Remove wholesale discount fields
		err = q.DisableWholesaleDiscount(ctx, sqlc.DisableWholesaleDiscountParams{
			Haswholesalediscount:  false,
			Wholesalediscounttype: sqlnull.String(""),
			Wholesalediscount:     sqlnull.Int64(0),
			UpdatedAt:             now,
			ID:                    input.VariantID,
			ProductID:             input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.wholesale_discount.removed event
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
			EventType:    "variant.wholesale_discount.removed",
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

	return &RemoveVariantWholesaleDiscountResult{
		VariantID: input.VariantID,
		Status:    "wholesale_discount_removed",
	}, nil
}
