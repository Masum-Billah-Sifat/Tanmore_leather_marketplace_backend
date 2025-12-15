// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/enable_variant_wholesale_mode_service.go
// 🧠 Enables wholesale mode on a product variant.
//     - Validates seller, product, variant, and category via snapshot
//     - Ensures wholesale mode is not already active
//     - Applies wholesale price and minimum quantity
//     - Emits variant.wholesale_mode.enabled event
//     - Returns variant_id and wholesale_enabled flag

package product_variant

import (
	"context"
	"encoding/json"
	"time"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_enable_wholesale_mode"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input from handler

type EnableWholesaleModeInput struct {
	UserID          uuid.UUID
	ProductID       uuid.UUID
	VariantID       uuid.UUID
	WholesalePrice  int64
	MinQtyWholesale int64
}

// ------------------------------------------------------------
// 👋 Result to return

type EnableWholesaleModeResult struct {
	VariantID        uuid.UUID
	WholesaleEnabled bool
	Status           string
}

// ------------------------------------------------------------
// 🔧 Dependencies

type EnableWholesaleModeServiceDeps struct {
	Repo repo.ProductVariantEnableWholesaleRepoInterface
}

// 🛠️ Service Struct
type EnableWholesaleModeService struct {
	Deps EnableWholesaleModeServiceDeps
}

// 🚀 Constructor
func NewEnableWholesaleModeService(deps EnableWholesaleModeServiceDeps) *EnableWholesaleModeService {
	return &EnableWholesaleModeService{Deps: deps}
}

// 🚀 Entrypoint
func (s *EnableWholesaleModeService) Start(
	ctx context.Context,
	input EnableWholesaleModeInput,
) (*EnableWholesaleModeResult, error) {

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

		// ✅ Must not already be wholesale enabled
		if snapshot.Haswholesaleenabled {
			return errors.NewConflictError("wholesale mode is already enabled for this variant")
		}

		// ------------------------------------------------------------
		// Step 3: Enable wholesale mode
		err = q.EnableWholesaleMode(ctx, sqlc.EnableWholesaleModeParams{
			WholesaleEnabled:      true,
			WholesalePrice:        sqlnull.Int64(input.WholesalePrice),
			MinQtyWholesale:       sqlnull.Int32(input.MinQtyWholesale),
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
		// Step 4: Emit event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":               input.UserID,
			"product_id":            input.ProductID,
			"variant_id":            input.VariantID,
			"wholesale_price":       input.WholesalePrice,
			"min_qty_wholesale":     input.MinQtyWholesale,
			"has_wholesale_enabled": true,
			// Pass optional discount fields even if empty
			"has_wholesale_discount":  false,
			"wholesale_discount_type": "",
			"wholesale_discount":      0,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.wholesale_mode.enabled",
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

	return &EnableWholesaleModeResult{
		VariantID:        input.VariantID,
		WholesaleEnabled: true,
		Status:           "wholesale_mode_enabled",
	}, nil
}
