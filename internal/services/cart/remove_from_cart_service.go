// ------------------------------------------------------------
// 📁 File: internal/services/cart/remove_cart_item_service.go
// 🧠 Handles removing a product variant from the customer's cart.
//     - Validates customer account (not banned, not archived)
//     - Fetches cart item by user and variant
//     - If not found or inactive, return error
//     - Else soft-deletes the row (sets quantity NULL + is_active = false)

package cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/remove_from_cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type RemoveCartItemInput struct {
	UserID    uuid.UUID
	VariantID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Result to return
type RemoveCartItemResult struct {
	VariantID uuid.UUID
	Status    string // "cart_item_removed"
}

// ------------------------------------------------------------
// 🧱 Dependencies
type RemoveCartItemServiceDeps struct {
	Repo repo.RemoveFromCartRepoInterface
}

// 🛠️ Service Definition
type RemoveCartItemService struct {
	Deps RemoveCartItemServiceDeps
}

// 🚀 Constructor
func NewRemoveCartItemService(deps RemoveCartItemServiceDeps) *RemoveCartItemService {
	return &RemoveCartItemService{Deps: deps}
}

// 🚀 Entrypoint
func (s *RemoveCartItemService) Start(
	ctx context.Context,
	input RemoveCartItemInput,
) (*RemoveCartItemResult, error) {

	now := timeutil.NowUTC()

	var result *RemoveCartItemResult

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
		// Step 2: Get cart item
		item, err := q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
			UserID:    input.UserID,
			VariantID: input.VariantID,
		})
		if err != nil {
			return errors.NewValidationError("cart_item", "item not found")
		}
		if !item.IsActive {
			return errors.NewConflictError("item already removed")
		}

		// ------------------------------------------------------------
		// Step 3: Deactivate cart item (soft delete)
		err = q.DeactivateCartItem(ctx, sqlc.DeactivateCartItemParams{
			UserID:           input.UserID,
			VariantID:        input.VariantID,
			UpdatedAt:        now,
			IsActive:         false,
			RequiredQuantity: sql.NullInt32{}, // this becomes NULL in SQL

		})
		if err != nil {
			return errors.NewTableError("cart_items.deactivate", err.Error())
		}

		result = &RemoveCartItemResult{
			VariantID: input.VariantID,
			Status:    "cart_item_removed",
		}
		return nil
	})

	// Return
	if err != nil {
		return nil, err
	}
	return result, nil
}
