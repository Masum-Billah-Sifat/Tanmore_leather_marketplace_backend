// 📁 File: internal/repository/cart/clear_cart/interface.go

package clear_cart

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type ClearCartRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	ClearCartItemsByUser(ctx context.Context, arg sqlc.ClearCartItemsByUserParams) error
}
