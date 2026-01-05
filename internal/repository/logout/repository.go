// ------------------------------------------------------------
// 📁 File: internal/repository/logout/logout_repository.go
// 🧠 This file provides the implementation of LogoutRepoInterface
//     using SQLC-generated methods for session termination.

package logout

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 LogoutRepository implements LogoutRepoInterface using sqlc
type LogoutRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor for LogoutRepository
func NewLogoutRepository(db *sql.DB) *LogoutRepository {
	return &LogoutRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *LogoutRepository) WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
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

// 🔎 Get user session by ID + user
func (r *LogoutRepository) GetSessionByIDAndUserID(ctx context.Context, sessionID, userID uuid.UUID) (sqlc.UserSession, error) {
	return r.q.GetSessionByIDAndUserID(ctx, sqlc.GetSessionByIDAndUserIDParams{
		ID:     sessionID,
		UserID: userID,
	})
}

// 🚫 Revoke session
func (r *LogoutRepository) RevokeUserSession(ctx context.Context, arg sqlc.RevokeUserSessionParams) error {
	return r.q.RevokeUserSession(ctx, arg)
}

// 🗑️ Deprecate refresh tokens
func (r *LogoutRepository) DeprecateRefreshTokensBySession(ctx context.Context, arg sqlc.DeprecateRefreshTokensBySessionParams) error {
	return r.q.DeprecateRefreshTokensBySession(ctx, arg)
}
