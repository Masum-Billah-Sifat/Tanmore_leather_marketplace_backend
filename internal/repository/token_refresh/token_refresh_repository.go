// ------------------------------------------------------------
// 📁 File: internal/repository/token_refresh/token_refresh_repository.go
// 🧠 This file provides the implementation of TokenRefreshRepoInterface
//     using SQLC-generated methods, aligned with Meta-grade standards.

package token_refresh

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 TokenRefreshRepository implements TokenRefreshRepoInterface using sqlc
type TokenRefreshRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor
func NewTokenRefreshRepository(db *sql.DB) *TokenRefreshRepository {
	return &TokenRefreshRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *TokenRefreshRepository) WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
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

// 🔍 Get refresh token by hashed value
func (r *TokenRefreshRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (sqlc.UserRefreshToken, error) {
	return r.q.GetRefreshTokenByHash(ctx, tokenHash)
}

// 👤 Get user by user ID
func (r *TokenRefreshRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

// 📲 Get session by session ID and user ID
func (r *TokenRefreshRepository) GetSessionByIDAndUserID(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (sqlc.UserSession, error) {
	return r.q.GetSessionByIDAndUserID(ctx, sqlc.GetSessionByIDAndUserIDParams{
		ID:     sessionID,
		UserID: userID,
	})
}

// 🚫 Deprecate old refresh token
func (r *TokenRefreshRepository) DeprecateRefreshTokenByID(ctx context.Context, arg sqlc.DeprecateRefreshTokenByIDParams) error {
	return r.q.DeprecateRefreshTokenByID(ctx, arg)
}

// 🔁 Insert new refresh token
func (r *TokenRefreshRepository) InsertRefreshToken(ctx context.Context, arg sqlc.InsertRefreshTokenParams) error {
	return r.q.InsertRefreshToken(ctx, arg)
}
