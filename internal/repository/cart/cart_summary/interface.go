package cart_summary

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// type CartSummaryRepoInterface interface {
// 	// 🔁 Transaction wrapper
// 	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

// 	// 🧑 Fetch customer user by ID
// 	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

// 	// 📦 Fetch active cart + snapshot for selected variants (supports guest)
// 	GetActiveCartVariantSnapshotsByOwnerAndVariantIDs(
// 		ctx context.Context,
// 		arg sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsParams,
// 	) ([]sqlc.GetActiveCartVariantSnapshotsByOwnerAndVariantIDsRow, error)
// }

type CartSummaryRepoInterface interface {
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

	GetActiveCartVariantSnapshotsByUserAndVariantIDs(
		ctx context.Context,
		arg sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsParams,
	) ([]sqlc.GetActiveCartVariantSnapshotsByUserAndVariantIDsRow, error)

	GetActiveCartVariantSnapshotsByGuestAndVariantIDs(
		ctx context.Context,
		arg sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsParams,
	) ([]sqlc.GetActiveCartVariantSnapshotsByGuestAndVariantIDsRow, error)
}
