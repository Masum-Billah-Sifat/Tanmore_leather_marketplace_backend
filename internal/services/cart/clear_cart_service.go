// ------------------------------------------------------------
// 📁 File: internal/services/cart/clear_cart_service.go
// 🧠 Handles clearing all active items from a customer's cart.
//     - Validates the user (not archived or banned)
//     - Soft-deletes (deactivates) all active cart rows for user
//     - Returns a result indicating cart cleared or already empty

package cart

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/clear_cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	// "database/sql" // <-- Ensure this is imported for sql.NullInt32

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type ClearCartInput struct {
	UserID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Result returned
type ClearCartResult struct {
	Status string // "cart_cleared"
}

// ------------------------------------------------------------
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

	// Default result
	result := &ClearCartResult{Status: "cart_cleared"}

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Validate user
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

		// Step 2: Clear cart items
		err = q.ClearCartItemsForUser(ctx, sqlc.ClearCartItemsForUserParams{
			UserID:    input.UserID,
			IsActive:  false,
			UpdatedAt: now,
			RequiredQuantity: sql.NullInt32{
				Valid: false, // sets to NULL
			},
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
