// ------------------------------------------------------------
// 📁 File: internal/services/cart/add_to_cart_service.go
// 🧠 Handles adding a product variant to the customer's cart.
//     - Validates customer account (not banned, not archived)
//     - Fetches variant snapshot (product + variant + seller)
//     - Validates moderation + availability + stock rules
//     - Reactivates old cart row if needed
//     - Otherwise inserts new cart row
//     - Returns status: added, already_in_cart, or reactivated

package cart

import (
	"context"
	"database/sql"
	"fmt"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/cart/add_to_cart"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// // 📥 Input from handler

type AddToCartInput struct {
	UserID           uuid.UUID
	ProductID        uuid.UUID
	VariantID        uuid.UUID
	RequiredQuantity int32
}

// ------------------------------------------------------------
// 📤 Result to return
type AddToCartResult struct {
	VariantID uuid.UUID
	Status    string // "added_to_cart", "already_in_cart", "cart_item_reactivated"
}

// ------------------------------------------------------------
// 🧱 Dependencies
type AddToCartServiceDeps struct {
	Repo repo.AddToCartRepoInterface
}

// 🛠️ Service Definition
type AddToCartService struct {
	Deps AddToCartServiceDeps
}

// 🚀 Constructor
func NewAddToCartService(deps AddToCartServiceDeps) *AddToCartService {
	return &AddToCartService{Deps: deps}
}

func (s *AddToCartService) Start(ctx context.Context, input AddToCartInput) (*AddToCartResult, error) {
	now := timeutil.NowUTC()
	var result *AddToCartResult

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {

		// Step 1️⃣: Validate user ONLY if authenticated
		// if input.UserID != nil {
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

		// Step 2️⃣: Fetch variant snapshot (UNCHANGED)
		snapshot, err := q.GetVariantSnapshotByProductIDAndVariantID(ctx, sqlc.GetVariantSnapshotByProductIDAndVariantIDParams{
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		// Step 3️⃣: Moderation + stock checks (UNCHANGED)
		if !snapshot.Issellerapproved || snapshot.Issellerarchived || snapshot.Issellerbanned {
			return errors.NewAuthError("seller moderation failed")
		}
		if !snapshot.Isproductapproved || snapshot.Isproductarchived || snapshot.Isproductbanned {
			return errors.NewValidationError("product", "product is not available for cart")
		}
		if snapshot.Isvariantarchived || !snapshot.Isvariantinstock {
			return errors.NewValidationError("variant", "variant is not in stock or archived")
		}
		if snapshot.Stockamount < input.RequiredQuantity {
			return errors.NewValidationError("required_quantity", "not enough stock available")
		}

		item, err := q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
			UserID:    input.UserID,
			VariantID: input.VariantID,
		})

		if err == nil {
			if item.IsActive {
				return errors.NewValidationError("variant", "item already exists in cart")
			}

			// Reactivate
			err := q.ReactivateCartItemByID(ctx, sqlc.ReactivateCartItemByIDParams{
				RequiredQuantity: sql.NullInt32{
					Int32: input.RequiredQuantity,
					Valid: true,
				},
				IsActive:  true,
				UpdatedAt: now,
				ID:        item.ID,
			})
			if err != nil {
				return errors.NewTableError("cart_items.reactivate", err.Error())
			}

			result = &AddToCartResult{
				VariantID: input.VariantID,
				Status:    "cart_item_reactivated",
			}
			return nil
		}

		err = q.InsertCartItem(ctx, sqlc.InsertCartItemParams{
			ID:        uuid.New(),
			UserID:    input.UserID,
			VariantID: input.VariantID,
			RequiredQuantity: sql.NullInt32{
				Int32: input.RequiredQuantity,
				Valid: true,
			},
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		})

		if err != nil {
			fmt.Println("❌ SQL Insert Error:", err)
			return errors.NewTableError("cart_items.insert", err.Error())
		}

		result = &AddToCartResult{
			VariantID: input.VariantID,
			Status:    "added_to_cart",
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
