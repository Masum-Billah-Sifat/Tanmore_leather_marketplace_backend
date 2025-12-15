// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/add_variant_retail_discount_service.go
// 🧠 Handles adding a retail discount to a product variant.
//     - Validates seller + product + variant + category via snapshot
//     - Applies the retail discount and type (flat or percentage)
//     - Emits variant.retail_discount.added event
//     - Returns variant_id, discount, and type

package product_variant

import (
	"context"
	"encoding/json"
	"time"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_add_discount"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input from handler

type AddVariantRetailDiscountInput struct {
	UserID             uuid.UUID
	ProductID          uuid.UUID
	VariantID          uuid.UUID
	RetailDiscount     int64
	RetailDiscountType string // "flat" or "percentage"
}

// ------------------------------------------------------------
// 👋 Result to return

type AddVariantRetailDiscountResult struct {
	VariantID          uuid.UUID
	RetailDiscount     int64
	RetailDiscountType string
	Status             string
}

// ------------------------------------------------------------
// 🔧 Dependencies

type AddVariantRetailDiscountServiceDeps struct {
	Repo repo.ProductVariantAddDiscountRepoInterface
}

// 🔧 Service Definition
type AddVariantRetailDiscountService struct {
	Deps AddVariantRetailDiscountServiceDeps
}

// 🚀 Constructor
func NewAddVariantRetailDiscountService(deps AddVariantRetailDiscountServiceDeps) *AddVariantRetailDiscountService {
	return &AddVariantRetailDiscountService{Deps: deps}
}

// 🚀 Entrypoint
func (s *AddVariantRetailDiscountService) Start(
	ctx context.Context,
	input AddVariantRetailDiscountInput,
) (*AddVariantRetailDiscountResult, error) {

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
		// Step 3: Update discount
		err = q.EnableRetailDiscount(ctx, sqlc.EnableRetailDiscountParams{
			Hasretaildiscount:  true,
			Retaildiscounttype: sqlnull.String(input.RetailDiscountType),
			Retaildiscount:     sqlnull.Int64(input.RetailDiscount),
			UpdatedAt:          now,
			ID:                 input.VariantID,
			ProductID:          input.ProductID,
		})
		if err != nil {
			return errors.NewTableError("product_variants.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.retail_discount.added event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":              input.UserID,
			"product_id":           input.ProductID,
			"variant_id":           input.VariantID,
			"retail_discount":      input.RetailDiscount,
			"retail_discount_type": input.RetailDiscountType,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.retail_discount.added",
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

	return &AddVariantRetailDiscountResult{
		VariantID:          input.VariantID,
		RetailDiscount:     input.RetailDiscount,
		RetailDiscountType: input.RetailDiscountType,
		Status:             "retail_discount_added",
	}, nil
}
