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
// 📥 Input from handler
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

		// Step 1️⃣: Validate user
		fmt.Println("🔍 Fetching user:", input.UserID)
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			fmt.Println("❌ User not found:", err)
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived {
			fmt.Println("❌ User is archived")
			return errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			fmt.Println("❌ User is banned")
			return errors.NewAuthError("user is banned")
		}

		// Step 2️⃣: Fetch variant snapshot
		fmt.Println("🔍 Fetching variant snapshot for product:", input.ProductID, "variant:", input.VariantID)
		snapshot, err := q.GetVariantSnapshotByProductIDAndVariantID(ctx, sqlc.GetVariantSnapshotByProductIDAndVariantIDParams{
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			fmt.Println("❌ Variant snapshot not found:", err)
			return errors.NewNotFoundError("variant snapshot")
		}

		// Step 3️⃣: Moderation checks
		if !snapshot.Issellerapproved || snapshot.Issellerarchived || snapshot.Issellerbanned {
			fmt.Println("❌ Seller moderation failed")
			return errors.NewAuthError("seller moderation failed")
		}
		if !snapshot.Isproductapproved || snapshot.Isproductarchived || snapshot.Isproductbanned {
			fmt.Println("❌ Product moderation failed")
			return errors.NewValidationError("product", "product is not available for cart")
		}
		if snapshot.Isvariantarchived || !snapshot.Isvariantinstock {
			fmt.Println("❌ Variant moderation failed")
			return errors.NewValidationError("variant", "variant is not in stock or archived")
		}
		if snapshot.Stockamount < input.RequiredQuantity {
			fmt.Println("❌ Insufficient stock. Available:", snapshot.Stockamount, "Requested:", input.RequiredQuantity)
			return errors.NewValidationError("required_quantity", "not enough stock available")
		}

		// Step 4️⃣: Check if cart item already exists
		fmt.Println("🔍 Checking if cart item exists")
		item, err := q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
			UserID:    input.UserID,
			VariantID: input.VariantID,
		})
		if err == nil {
			fmt.Println("ℹ️ Cart item already exists:", item.ID, "isActive:", item.IsActive)
			if item.IsActive {
				return errors.NewValidationError("variant", "item already exists in cart")
			}

			// Reactivate soft-deleted cart item
			fmt.Println("♻️ Reactivating cart item")
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
				fmt.Println("❌ Failed to reactivate cart item:", err)
				return errors.NewTableError("cart_items.reactivate", err.Error())
			}

			result = &AddToCartResult{
				VariantID: input.VariantID,
				Status:    "cart_item_reactivated",
			}
			return nil
		} else {
			fmt.Println("ℹ️ No existing cart item found. Proceeding to insert.")
		}

		// Step 5️⃣: Insert new cart item
		fmt.Println("🛒 Inserting new cart item")
		_, err = q.InsertCartItem(ctx, sqlc.InsertCartItemParams{
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
			fmt.Println("❌ Failed to insert cart item:", err)
			return errors.NewTableError("cart_items.insert", err.Error())
		}

		result = &AddToCartResult{
			VariantID: input.VariantID,
			Status:    "added_to_cart",
		}
		return nil
	})

	// Final step: return result or error
	if err != nil {
		fmt.Println("❌ Service error:", err)
		return nil, err
	}
	fmt.Println("✅ Service succeeded:", result)
	return result, nil
}

// // 🚀 Entrypoint
// func (s *AddToCartService) Start(
// 	ctx context.Context,
// 	input AddToCartInput,
// ) (*AddToCartResult, error) {

// 	now := timeutil.NowUTC()

// 	var result *AddToCartResult

// 	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
// 		// ------------------------------------------------------------
// 		// Step 1: Validate customer
// 		user, err := q.GetUserByID(ctx, input.UserID)
// 		if err != nil {
// 			return errors.NewNotFoundError("user")
// 		}
// 		if user.IsArchived {
// 			return errors.NewAuthError("user is archived")
// 		}
// 		if user.IsBanned {
// 			return errors.NewAuthError("user is banned")
// 		}

// 		// ------------------------------------------------------------
// 		// Step 2: Fetch snapshot for product + variant + seller
// 		snapshot, err := q.GetVariantSnapshotByProductIDAndVariantID(ctx, sqlc.GetVariantSnapshotByProductIDAndVariantIDParams{
// 			Productid: input.ProductID,
// 			Variantid: input.VariantID,
// 		})
// 		if err != nil {
// 			return errors.NewNotFoundError("variant snapshot")
// 		}

// 		// Moderation checks
// 		if !snapshot.Issellerapproved || snapshot.Issellerarchived || snapshot.Issellerbanned {
// 			return errors.NewAuthError("seller moderation failed")
// 		}
// 		if !snapshot.Isproductapproved || snapshot.Isproductarchived || snapshot.Isproductbanned {
// 			return errors.NewValidationError("product", "product is not available for cart")
// 		}
// 		if snapshot.Isvariantarchived || !snapshot.Isvariantinstock {
// 			return errors.NewValidationError("variant", "variant is not in stock or archived")
// 		}
// 		if snapshot.Stockamount < input.RequiredQuantity {
// 			return errors.NewValidationError("required_quantity", "not enough stock available")
// 		}

// 		// ------------------------------------------------------------
// 		// Step 3: Try to find existing cart item
// 		item, err := q.GetCartItemByUserAndVariant(ctx, sqlc.GetCartItemByUserAndVariantParams{
// 			UserID:    input.UserID,
// 			VariantID: input.VariantID,
// 		})

// 		if err == nil {
// 			if item.IsActive {
// 				return errors.NewValidationError(
// 					"variant",
// 					"item already exists in cart",
// 				)
// 			}

// 			// Case B: Reactivate
// 			err := q.ReactivateCartItemByID(ctx, sqlc.ReactivateCartItemByIDParams{
// 				RequiredQuantity: sql.NullInt32{
// 					Int32: input.RequiredQuantity,
// 					Valid: true,
// 				},
// 				IsActive:  true,
// 				UpdatedAt: now,
// 				ID:        item.ID,
// 			})
// 			if err != nil {
// 				return errors.NewTableError("cart_items.reactivate", err.Error())
// 			}

// 			result = &AddToCartResult{
// 				VariantID: input.VariantID,
// 				Status:    "cart_item_reactivated",
// 			}
// 			return nil
// 		}

// 		// ------------------------------------------------------------
// 		// Step 4: Insert new cart item
// 		_, err = q.InsertCartItem(ctx, sqlc.InsertCartItemParams{
// 			UserID:    input.UserID,
// 			VariantID: input.VariantID,
// 			RequiredQuantity: sql.NullInt32{
// 				Int32: input.RequiredQuantity,
// 				Valid: true,
// 			},
// 			IsActive:  true,
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		})
// 		if err != nil {
// 			return errors.NewTableError("cart_items.insert", err.Error())
// 		}

// 		result = &AddToCartResult{
// 			VariantID: input.VariantID,
// 			Status:    "added_to_cart",
// 		}
// 		return nil
// 	})

// 	// Return result or error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
