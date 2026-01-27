// ------------------------------------------------------------
// 📁 File: internal/services/cart/update_cart_quantity_service.go
// 🧠 Handles updating quantity of an existing cart item.

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/update_to_cart"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateCartQuantityInput struct {
	UserID           *uuid.UUID
	GuestUserID      *uuid.UUID
	VariantID        uuid.UUID
	RequiredQuantity int32
}

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

func (s *UpdateCartQuantityService) Start(
	ctx context.Context,
	input UpdateCartQuantityInput,
) (*UpdateCartQuantityResult, error) {

	now := timeutil.NowUTC()
	var result *UpdateCartQuantityResult

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate user if logged in
		if input.UserID != nil {
			user, err := q.GetUserByID(ctx, *input.UserID)
			if err != nil {
				return errors.NewNotFoundError("user")
			}
			if user.IsArchived {
				return errors.NewAuthError("user is archived")
			}
			if user.IsBanned {
				return errors.NewAuthError("user is banned")
			}
		}

		// ------------------------------------------------------------
		// Step 2: Fetch variant snapshot
		snapshot, err := q.GetVariantSnapshotByVariantID(ctx, input.VariantID)
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

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
		// Step 3: Get cart item by user or guest
		var item sqlc.GetCartItemByUserAndVariantRow

		if input.UserID != nil {
			item, err = q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
				UserID:    sqlnull.UUIDPtr(input.UserID),
				VariantID: input.VariantID,
			})
		} else {
			guestItem, gErr := q.GetCartItemByGuestAndVariant(ctx, sqlc.GetCartItemByGuestAndVariantParams{
				GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
				VariantID:   input.VariantID,
			})
			item = sqlc.GetCartItemByUserAndVariantRow(guestItem)
			err = gErr
		}

		if err != nil {
			return errors.NewNotFoundError("item_not_found_or_archived")
		}
		if !item.IsActive {
			return errors.NewValidationError("cart", "cart item is archived")
		}

		// ------------------------------------------------------------
		// Step 4: Update quantity based on owner type
		if input.UserID != nil {
			err = q.UpdateCartQuantityForUser(ctx, sqlc.UpdateCartQuantityForUserParams{
				UserID:    sqlnull.UUIDPtr(input.UserID),
				VariantID: input.VariantID,
				// RequiredQuantity: input.RequiredQuantity,
				RequiredQuantity: sqlnull.Int32(int64(input.RequiredQuantity)),
				UpdatedAt:        now,
			})
		} else {
			err = q.UpdateCartQuantityForGuest(ctx, sqlc.UpdateCartQuantityForGuestParams{
				GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
				VariantID:   input.VariantID,
				// RequiredQuantity: input.RequiredQuantity,
				RequiredQuantity: sqlnull.Int32(int64(input.RequiredQuantity)),
				UpdatedAt:        now,
			})
		}

		if err != nil {
			return errors.NewTableError("cart_items.update_quantity", err.Error())
		}

		result = &UpdateCartQuantityResult{
			VariantID:       input.VariantID,
			UpdatedQuantity: input.RequiredQuantity,
			Status:          "cart_item_updated",
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
