// 📁 File: internal/services/cart/clear_cart_service.go

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/clear_cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

type ClearCartInput struct {
	UserID uuid.UUID
}

type ClearCartResult struct {
	Status string // "cart_cleared"
}

type ClearCartServiceDeps struct {
	Repo repo.ClearCartRepoInterface
}

type ClearCartService struct {
	Deps ClearCartServiceDeps
}

func NewClearCartService(deps ClearCartServiceDeps) *ClearCartService {
	return &ClearCartService{Deps: deps}
}

func (s *ClearCartService) Start(
	ctx context.Context,
	input ClearCartInput,
) (*ClearCartResult, error) {

	now := timeutil.NowUTC()
	result := &ClearCartResult{
		Status: "cart_cleared",
	}

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

		// 2️⃣ Clear cart
		err = q.ClearCartItemsByUser(ctx, sqlc.ClearCartItemsByUserParams{
			UserID:    input.UserID,
			UpdatedAt: now,
		})
		if err != nil {
			return errors.NewTableError("cart_items.clear", err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
