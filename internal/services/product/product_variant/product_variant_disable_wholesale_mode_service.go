// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/disable_variant_wholesale_mode_service.go
// 🧠 Disables wholesale mode on a product variant.
//     - Validates seller, product, variant, and category via snapshot
//     - Ensures wholesale mode is currently enabled
//     - Resets wholesale-related fields to NULL/default
//     - Emits variant.wholesale_mode.disabled event
//     - Returns variant_id and wholesale_enabled flag

package product_variant

import (
	"context"
	"encoding/json"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_disable_wholesale"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📅 Input from handler

type DisableWholesaleModeInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	VariantID uuid.UUID
}

// ------------------------------------------------------------
// 👋 Result to return

type DisableWholesaleModeResult struct {
	VariantID        uuid.UUID
	WholesaleEnabled bool
	Status           string
}

// ------------------------------------------------------------
// 🔧 Dependencies

type DisableWholesaleModeServiceDeps struct {
	Repo repo.ProductVariantDisableWholesaleRepoInterface
}

// 🛠️ Service Struct
type DisableWholesaleModeService struct {
	Deps DisableWholesaleModeServiceDeps
}

// 🚀 Constructor
func NewDisableWholesaleModeService(deps DisableWholesaleModeServiceDeps) *DisableWholesaleModeService {
	return &DisableWholesaleModeService{Deps: deps}
}

// 🚀 Entrypoint
func (s *DisableWholesaleModeService) Start(
	ctx context.Context,
	input DisableWholesaleModeInput,
) (*DisableWholesaleModeResult, error) {

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
		// Step 2: Moderation and status checks
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

		if !snapshot.Haswholesaleenabled {
			return errors.NewConflictError("wholesale mode is already disabled for this variant")
		}

		// ------------------------------------------------------------
		// Step 3: Disable wholesale mode and clear all related fields
		err = q.DisableWholesaleMode(ctx, sqlc.DisableWholesaleModeParams{
			WholesaleEnabled:      false,
			WholesalePrice:        sqlnull.Int64(0),
			MinQtyWholesale:       sqlnull.Int32(0),
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
		// Step 4: Emit event with full payload (for background processors)
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
			EventType:    "variant.wholesale_mode.disabled",
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

	return &DisableWholesaleModeResult{
		VariantID:        input.VariantID,
		WholesaleEnabled: false,
		Status:           "wholesale_mode_disabled",
	}, nil
}
