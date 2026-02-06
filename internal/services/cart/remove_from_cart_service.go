package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/remove_from_cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

type RemoveCartItemInput struct {
	UserID    uuid.UUID
	VariantID uuid.UUID
}

type RemoveCartItemResult struct {
	VariantID uuid.UUID
	Status    string
}

type RemoveCartItemServiceDeps struct {
	Repo repo.RemoveFromCartRepoInterface
}

type RemoveCartItemService struct {
	Deps RemoveCartItemServiceDeps
}

func NewRemoveCartItemService(deps RemoveCartItemServiceDeps) *RemoveCartItemService {
	return &RemoveCartItemService{Deps: deps}
}

func (s *RemoveCartItemService) Start(
	ctx context.Context,
	input RemoveCartItemInput,
) (*RemoveCartItemResult, error) {
	now := timeutil.NowUTC()

	var result *RemoveCartItemResult

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// 🔒 Validate user
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.ErrAuthUserNotFound()
		}
		if user.IsArchived {
			return errors.ErrAuthArchivedUser()
		}
		if user.IsBanned {
			return errors.ErrAuthBannedUser()
		}

		// 🛒 Check cart item
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

		// ❌ Deactivate
		err = q.DeactivateCartItemByUserAndVariant(ctx, sqlc.DeactivateCartItemByUserAndVariantParams{
			UserID:    input.UserID,
			VariantID: input.VariantID,
			UpdatedAt: now,
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

	if err != nil {
		return nil, err
	}
	return result, nil
}
