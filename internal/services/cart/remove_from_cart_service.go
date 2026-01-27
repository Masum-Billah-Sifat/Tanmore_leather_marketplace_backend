package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/remove_from_cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// 📥 Input from handler
type RemoveCartItemInput struct {
	UserID      *uuid.UUID
	GuestUserID *uuid.UUID
	VariantID   uuid.UUID
}

// 📤 Result to return
type RemoveCartItemResult struct {
	VariantID uuid.UUID
	Status    string // "cart_item_removed"
}

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
		// Step 1: Validate user ONLY if authenticated
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

		// Step 2: Get cart item (split between user/guest)
		var item sqlc.GetCartItemByUserAndVariantRow
		var err error

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
			return errors.NewValidationError("cart_item", "item not found")
		}
		if !item.IsActive {
			return errors.NewConflictError("item already removed")
		}

		// Step 3: Deactivate
		if input.UserID != nil {
			err = q.DeactivateCartItemByUserAndVariant(ctx, sqlc.DeactivateCartItemByUserAndVariantParams{
				UserID:    sqlnull.UUIDPtr(input.UserID),
				VariantID: input.VariantID,
				UpdatedAt: now,
			})
		} else {
			err = q.DeactivateCartItemByGuestAndVariant(ctx, sqlc.DeactivateCartItemByGuestAndVariantParams{
				GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
				VariantID:   input.VariantID,
				UpdatedAt:   now,
			})
		}

		if err != nil {
			return errors.NewTableError("cart_items.deactivate", err.Error())
		}

		result = &RemoveCartItemResult{
			VariantID: input.VariantID,
			Status:    "cart_item_removed",
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
