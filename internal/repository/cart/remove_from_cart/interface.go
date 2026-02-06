package remove_from_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// internal/repository/cart/remove_from_cart/interface.go

type RemoveFromCartRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	GetCartItemByUserAndVariant(ctx context.Context, userID uuid.UUID, variantID uuid.UUID) (sqlc.CartItem, error)
	DeactivateCartItemByUserAndVariant(ctx context.Context, arg sqlc.DeactivateCartItemByUserAndVariantParams) error
}
