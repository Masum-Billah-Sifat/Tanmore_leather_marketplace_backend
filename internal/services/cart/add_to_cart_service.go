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
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// // 📥 Input from handler
// type AddToCartInput struct {
// 	UserID           uuid.UUID
// 	ProductID        uuid.UUID
// 	VariantID        uuid.UUID
// 	RequiredQuantity int32
// }

// 📥 Input from handler
type AddToCartInput struct {
	UserID      *uuid.UUID
	GuestUserID *uuid.UUID

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

		// Step 4️⃣: Check existing cart item (USER OR GUEST)
		// Step 4️⃣: Check existing cart item (split for user vs guest)
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
			item = sqlc.GetCartItemByUserAndVariantRow(guestItem) // type-compatible (same fields)
			err = gErr
		}

		if err == nil {
			if item.IsActive {
				return errors.NewValidationError("variant", "item already exists in cart")
			}

			// Reactivate
			err := q.ReactivateCartItemByID(ctx, sqlc.ReactivateCartItemByIDParams{
				ID: item.ID,
				RequiredQuantity: sql.NullInt32{
					Int32: input.RequiredQuantity,
					Valid: true,
				},
				UpdatedAt: now,
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

		fmt.Println("DEBUG – UserID:", input.UserID)
		fmt.Println("DEBUG – GuestUserID:", input.GuestUserID)

		// Step 5️⃣: Insert new cart item
		_, err = q.InsertCartItem(ctx, sqlc.InsertCartItemParams{
			UserID:      sqlnull.UUIDPtr(input.UserID),
			GuestUserID: sqlnull.UUIDPtr(input.GuestUserID),
			VariantID:   input.VariantID,
			RequiredQuantity: sql.NullInt32{
				Int32: input.RequiredQuantity,
				Valid: true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		})
		// if err != nil {
		// 	return errors.NewTableError("cart_items.insert", err.Error())
		// }

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
