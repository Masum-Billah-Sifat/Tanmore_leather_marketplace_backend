// ------------------------------------------------------------
// 📁 File: internal/services/cart/clear_cart_service.go
// 🧠 Handles clearing all active items from a customer's cart.
//     - Validates the user (if logged in)
//     - Soft-deletes (deactivates) all active cart rows for user or guest
//     - Returns a result indicating cart cleared

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/clear_cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// 📥 Input from handler
type ClearCartInput struct {
	UserID      *uuid.UUID
	GuestUserID *uuid.UUID
}

// 📤 Result returned
type ClearCartResult struct {
	Status string // "cart_cleared"
}

// 🧱 Dependencies
type ClearCartServiceDeps struct {
	Repo repo.ClearCartRepoInterface
}

// 🛠️ Service Definition
type ClearCartService struct {
	Deps ClearCartServiceDeps
}

// 🚀 Constructor
func NewClearCartService(deps ClearCartServiceDeps) *ClearCartService {
	return &ClearCartService{Deps: deps}
}

// 🚀 Entrypoint
func (s *ClearCartService) Start(
	ctx context.Context,
	input ClearCartInput,
) (*ClearCartResult, error) {

	now := timeutil.NowUTC()
	result := &ClearCartResult{
		Status: "cart_cleared",
	}

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1️⃣: Validate user ONLY if logged in
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

		// Step 2️⃣: Clear cart based on user or guest
		var err error
		if input.UserID != nil {
			err = q.ClearCartItemsByUser(ctx, sqlc.ClearCartItemsByUserParams{
				UserID:    sqlnull.UUIDPtr(input.UserID),
				UpdatedAt: now,
			})
		} else {
			err = q.ClearCartItemsByGuest(ctx, sqlc.ClearCartItemsByGuestParams{
				GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
				UpdatedAt:   now,
			})
		}

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
