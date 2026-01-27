// ------------------------------------------------------------
// 📁 File: internal/services/cart/merge_guest_cart_service.go
// 🧠 Handles merging guest cart items into authenticated user's cart.
//     - Validates customer (not banned, not archived)
//     - Fetches guest cart items
//     - Inserts guest items for authenticated user
//     - Deprecates old guest items
//     - Returns merge summary

package cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/merge_guest_cart"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// 📥 Input

type MergeGuestCartInput struct {
	UserID      uuid.UUID
	GuestUserID uuid.UUID
}

// 📤 Result

type MergeGuestCartResult struct {
	Status      string `json:"status"`
	MergedCount int    `json:"merged_count,omitempty"`
}

// ------------------------------------------------------------
// 🧱 Dependencies
type MergeGuestCartServiceDeps struct {
	Repo repo.MergeGuestCartRepoInterface
}

// 🛠️ Service Definition
type MergeGuestCartService struct {
	Deps MergeGuestCartServiceDeps
}

// 🚀 Constructor
func NewMergeGuestCartService(deps MergeGuestCartServiceDeps) *MergeGuestCartService {
	return &MergeGuestCartService{Deps: deps}
}

// 🔁 Entry Point
func (s *MergeGuestCartService) Start(
	ctx context.Context,
	input MergeGuestCartInput,
) (*MergeGuestCartResult, error) {
	now := timeutil.NowUTC()

	var mergedCount int

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1️⃣: Validate customer moderation
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

		// Step 2️⃣: Fetch guest cart items
		items, err := q.GetActiveGuestCartItems(ctx, sqlnull.UUID(input.GuestUserID))
		if err != nil {
			return errors.NewServerError("could not read guest cart")
		}
		if len(items) == 0 {
			return errors.NewServerError("no guest cart items to merge")
		}

		// Step 3️⃣: Insert each guest item for logged-in user
		for _, item := range items {
			_, err := q.InsertCartItem(ctx, sqlc.InsertCartItemParams{
				UserID:           sqlnull.UUIDPtr(&input.UserID),
				GuestUserID:      sqlnull.UUIDPtr(&input.GuestUserID),
				VariantID:        item.VariantID,
				RequiredQuantity: item.RequiredQuantity,
				// IsActive:         true,
				// IsDeprecated:     false,
				CreatedAt: now,
				UpdatedAt: now,
			})
			if err != nil {
				return errors.NewTableError("insert_cart_item", err.Error())
			}
			mergedCount++
		}

		// Step 4️⃣: Deprecate old guest items
		err = q.DeprecateGuestCartItems(ctx, sqlc.DeprecateGuestCartItemsParams{
			GuestUserID:  sqlnull.UUID(input.GuestUserID),
			IsActive:     true,
			IsDeprecated: true,
			UpdatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("deprecate_guest_items", err.Error())
		}

		return nil
	})

	if err != nil {
		if errors.IsCustomCode(err, "no_items_to_merge") {
			return &MergeGuestCartResult{Status: "no_items_to_merge"}, nil
		}
		return nil, err
	}

	return &MergeGuestCartResult{
		Status:      "guest_cart_merged",
		MergedCount: mergedCount,
	}, nil
}
