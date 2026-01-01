// ------------------------------------------------------------
// 📁 File: internal/services/cart/update_cart_quantity_service.go
// 🧠 Handles updating quantity of an existing cart item.
//     - Validates customer (not banned/archived)
//     - Fetches variant snapshot by variant_id
//     - Validates moderation and availability rules
//     - Verifies existing active cart row
//     - Updates quantity + timestamp

package cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/update_to_cart"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateCartQuantityInput struct {
	UserID           uuid.UUID
	VariantID        uuid.UUID
	RequiredQuantity int32
}

// ------------------------------------------------------------
// 📤 Result to return
type UpdateCartQuantityResult struct {
	VariantID       uuid.UUID
	UpdatedQuantity int32
	Status          string // "cart_item_updated"
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateCartQuantityServiceDeps struct {
	Repo repo.UpdateCartQuantityRepoInterface
}

// 🛠️ Service Definition
type UpdateCartQuantityService struct {
	Deps UpdateCartQuantityServiceDeps
}

// 🚀 Constructor
func NewUpdateCartQuantityService(deps UpdateCartQuantityServiceDeps) *UpdateCartQuantityService {
	return &UpdateCartQuantityService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateCartQuantityService) Start(
	ctx context.Context,
	input UpdateCartQuantityInput,
) (*UpdateCartQuantityResult, error) {

	now := timeutil.NowUTC()

	var result *UpdateCartQuantityResult

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate customer
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived {
			return errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			return errors.NewAuthError("user is banned")
		}

		// ------------------------------------------------------------
		// Step 2: Fetch variant snapshot
		snapshot, err := q.GetVariantSnapshotByVariantID(ctx, input.VariantID)
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		// Moderation checks
		if !snapshot.Issellerapproved || snapshot.Issellerarchived || snapshot.Issellerbanned {
			return errors.NewAuthError("seller moderation failed")
		}
		if !snapshot.Isproductapproved || snapshot.Isproductarchived || snapshot.Isproductbanned {
			return errors.NewValidationError("product", "product is not available for update")
		}
		if snapshot.Isvariantarchived || !snapshot.Isvariantinstock {
			return errors.NewValidationError("variant", "variant is not in stock or archived")
		}
		if snapshot.Stockamount < input.RequiredQuantity {
			return errors.NewValidationError("required_quantity", "not enough stock available")
		}

		// ------------------------------------------------------------
		// Step 3: Validate cart row exists and is active
		item, err := q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
			UserID:    input.UserID,
			VariantID: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("item_not_found_or_archived")
		}
		if !item.IsActive {
			return errors.NewValidationError("cart", "cart item is archived")
		}

		// ------------------------------------------------------------
		// Step 4: Update quantity
		err = q.UpdateCartQuantity(ctx, sqlc.UpdateCartQuantityParams{
			RequiredQuantity: sql.NullInt32{
				Int32: input.RequiredQuantity,
				Valid: true,
			},
			UpdatedAt: now,
			UserID:    input.UserID,
			VariantID: input.VariantID,
			IsActive:  true,
		})
		if err != nil {
			return errors.NewTableError("cart_items.update_quantity", err.Error())
		}

		// Return result
		result = &UpdateCartQuantityResult{
			VariantID:       input.VariantID,
			UpdatedQuantity: input.RequiredQuantity,
			Status:          "cart_item_updated",
		}
		return nil
	})

	// Return result or error
	if err != nil {
		return nil, err
	}

	return result, nil
}
