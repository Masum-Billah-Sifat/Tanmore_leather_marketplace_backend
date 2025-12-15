// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/add_variant_wholesale_discount_service.go
// 🧠 Adds wholesale discount to a product variant.
//     - Validates seller, product, variant, and category via snapshot
//     - Ensures wholesale mode is already enabled
//     - Applies wholesale discount value and type
//     - Emits variant.wholesale_discount.added event
//     - Returns variant_id and discount metadata

package product_variant

import (
	"context"
	"encoding/json"
	"time"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_add_wholesale_discount"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input from handler

type AddWholesaleDiscountInput struct {
	UserID                uuid.UUID
	ProductID             uuid.UUID
	VariantID             uuid.UUID
	WholesaleDiscount     int64
	WholesaleDiscountType string // "percentage" or "flat"
}

// ------------------------------------------------------------
// 👋 Result to return

type AddWholesaleDiscountResult struct {
	VariantID             uuid.UUID
	WholesaleDiscount     int64
	WholesaleDiscountType string
	Status                string
}

// ------------------------------------------------------------
// 🔧 Dependencies

type AddWholesaleDiscountServiceDeps struct {
	Repo repo.ProductVariantAddWholesaleDiscountRepoInterface
}

// 🛠️ Service Struct
type AddWholesaleDiscountService struct {
	Deps AddWholesaleDiscountServiceDeps
}

// 🚀 Constructor
func NewAddWholesaleDiscountService(deps AddWholesaleDiscountServiceDeps) *AddWholesaleDiscountService {
	return &AddWholesaleDiscountService{Deps: deps}
}

// 🚀 Entrypoint
func (s *AddWholesaleDiscountService) Start(
	ctx context.Context,
	input AddWholesaleDiscountInput,
) (*AddWholesaleDiscountResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Fetch variant snapshot
		snapshot, err := q.GetVariantSnapshot(ctx, sqlc.GetVariantSnapshotParams{
			Sellerid:  input.UserID,
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		// ------------------------------------------------------------
		// Step 2: Moderation checks
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

		// ✅ Must already have wholesale mode enabled
		if !snapshot.Haswholesaleenabled {
			return errors.NewConflictError("wholesale mode must be enabled before adding a discount")
		}

		// ------------------------------------------------------------
		// Step 3: Apply wholesale discount
		err = q.EnableWholesaleDiscount(ctx, sqlc.EnableWholesaleDiscountParams{
			Haswholesalediscount:  true,
			Wholesalediscounttype: sqlnull.String(input.WholesaleDiscountType),
			Wholesalediscount:     sqlnull.Int64(input.WholesaleDiscount),
			UpdatedAt:             now,
			ID:                    input.VariantID,
			ProductID:             input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":                 input.UserID,
			"product_id":              input.ProductID,
			"variant_id":              input.VariantID,
			"wholesale_discount":      input.WholesaleDiscount,
			"wholesale_discount_type": input.WholesaleDiscountType,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.wholesale_discount.added",
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

	return &AddWholesaleDiscountResult{
		VariantID:             input.VariantID,
		WholesaleDiscount:     input.WholesaleDiscount,
		WholesaleDiscountType: input.WholesaleDiscountType,
		Status:                "wholesale_discount_added",
	}, nil
}
