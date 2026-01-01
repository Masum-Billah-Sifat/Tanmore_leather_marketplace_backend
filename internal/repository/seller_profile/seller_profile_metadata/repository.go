// ------------------------------------------------------------
// 📁 File: internal/repository/seller_profile/seller_profile_metadata/seller_profile_repository.go
// 🧠 Concrete implementation of SellerProfileMetadataRepoInterface
//     Uses SQLC to manage user fetch, profile insert, status update, and event logging.

package seller_profile_metadata

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 SellerProfileMetadataRepository implements SellerProfileMetadataRepoInterface
type SellerProfileMetadataRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewSellerProfileMetadataRepository(db *sql.DB) *SellerProfileMetadataRepository {
	return &SellerProfileMetadataRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction wrapper
func (r *SellerProfileMetadataRepository) WithTx(
	ctx context.Context,
	fn func(q *sqlc.Queries) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := sqlc.New(tx)

	if err := fn(qtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// 🔍 Fetch user by ID
func (r *SellerProfileMetadataRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// 📝 Insert seller profile metadata
func (r *SellerProfileMetadataRepository) InsertSellerProfileMetadata(
	ctx context.Context,
	arg sqlc.InsertSellerProfileMetadataParams,
) (uuid.UUID, error) {
	return r.q.InsertSellerProfileMetadata(ctx, arg)
}

// ✅ Mark user as profile created
func (r *SellerProfileMetadataRepository) UpdateSellerProfileCreated(
	ctx context.Context,
	arg sqlc.UpdateSellerProfileCreatedParams,
) error {
	return r.q.UpdateSellerProfileCreated(ctx, arg)
}
