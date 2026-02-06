package update_cart_quantity

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type UpdateCartQuantityRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	GetVariantSnapshotByVariantID(ctx context.Context, variantID uuid.UUID) (sqlc.ProductVariantSnapshot, error)

	GetCartItemByUserAndVariant(
		ctx context.Context,
		userID uuid.UUID,
		variantID uuid.UUID,
	) (sqlc.CartItem, error)

	UpdateCartQuantityForUser(
		ctx context.Context,
		arg sqlc.UpdateCartQuantityForUserParams,
	) error
}
