// ------------------------------------------------------------
// 📁 File: internal/repository/product/product_add_media/interface.go
// 🧠 Repository interface for adding media (image or promo video) to a product.
//     Handles user/product validation, promo video existence check, insert, and event.

package product_add_media

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ProductAddMediaRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Fetch seller by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	// 📦 Fetch product by ID and seller ID
	GetProductByIDAndSellerID(
		ctx context.Context,
		arg sqlc.GetProductByIDAndSellerIDParams,
	) (sqlc.Product, error)

	// 📽️ Check if a promo video already exists for this product
	GetPromoVideoByProductID(
		ctx context.Context,
		arg sqlc.GetPromoVideoByProductIDParams,
	) (uuid.UUID, error)

	// 📥 Insert into product_medias table
	InsertProductMedia(
		ctx context.Context,
		arg sqlc.InsertProductMediaParams,
	) (uuid.UUID, error)

	// 📨 Insert event into outbox
	InsertEvent(ctx context.Context, arg sqlc.InsertEventParams) error
}
