// ------------------------------------------------------------
// 📁 File: internal/repository/seller_profile/seller_profile_metadata/interface.go
// 🧠 Repository interface for creating seller profile metadata.
//     Includes user fetch, seller profile insert, user update, and event insertion.

package seller_profile_metadata

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type SellerProfileMetadataRepoInterface interface {
	// 🔁 Transaction wrapper
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🔍 Fetch user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)

	InsertSellerProfileMetadata(
		ctx context.Context,
		arg sqlc.InsertSellerProfileMetadataParams,
	) (uuid.UUID, error)

	// ✅ Mark user's profile as created
	UpdateSellerProfileCreated(
		ctx context.Context,
		arg sqlc.UpdateSellerProfileCreatedParams,
	) error
}
