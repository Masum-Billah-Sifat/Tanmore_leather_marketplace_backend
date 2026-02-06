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

type UpdateCartQuantityInput struct {
	UserID           uuid.UUID
	VariantID        uuid.UUID
	RequiredQuantity int32
}

type UpdateCartQuantityResult struct {
	VariantID       uuid.UUID
	UpdatedQuantity int32
	Status          string
}

type UpdateCartQuantityServiceDeps struct {
	Repo repo.UpdateCartQuantityRepoInterface
}

type UpdateCartQuantityService struct {
	Deps UpdateCartQuantityServiceDeps
}

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

		// 1️⃣ Validate user
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

		// 2️⃣ Validate variant snapshot
		snapshot, err := q.GetVariantSnapshotByVariantID(ctx, input.VariantID)
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		if snapshot.Isvariantarchived || !snapshot.Isvariantinstock {
			return errors.NewValidationError("variant", "variant unavailable")
		}
		if snapshot.Stockamount < input.RequiredQuantity {
			return errors.NewValidationError("required_quantity", "not enough stock")
		}

		// 3️⃣ Fetch cart item (USER ONLY)
		item, err := q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
			UserID:    input.UserID,
			VariantID: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("item_not_found")
		}

		if !item.IsActive {
			return errors.NewValidationError("cart", "cart item is inactive")
		}

		// 4️⃣ Update quantity
		err = q.UpdateCartQuantityForUser(ctx, sqlc.UpdateCartQuantityForUserParams{
			UserID:           input.UserID,
			VariantID:        input.VariantID,
			RequiredQuantity: sqlnull.Int32(int64(input.RequiredQuantity)),
			UpdatedAt:        now,
		})
		if err != nil {
			return errors.NewTableError("cart_items.update", err.Error())
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
