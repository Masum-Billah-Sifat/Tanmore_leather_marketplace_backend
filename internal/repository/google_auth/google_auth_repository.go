// ------------------------------------------------------------
// 📁 File: internal/repository/google_auth/google_auth_repository.go
// 🧠 This file provides the implementation of GoogleAuthRepoInterface
//     using SQLC-generated methods, aligned with Meta-grade standards.

package google_auth

import (
	"context"
	"database/sql"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

// 📦 GoogleAuthRepository implements GoogleAuthRepoInterface using sqlc
type GoogleAuthRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

// 🚀 Constructor for GoogleAuthRepository
func NewGoogleAuthRepository(db *sql.DB) *GoogleAuthRepository {
	return &GoogleAuthRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// 🔁 Transaction handler
func (r *GoogleAuthRepository) WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
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

// 👤 Fetch existing user by Google ID
func (r *GoogleAuthRepository) GetUserByGoogleID(ctx context.Context, googleID string) (sqlc.User, error) {
	return r.q.GetUserByGoogleID(ctx, googleID)
}

// ➕ Insert new user
func (r *GoogleAuthRepository) InsertUser(ctx context.Context, arg sqlc.InsertUserParams) (uuid.UUID, error) {
	return r.q.InsertUser(ctx, arg)
}

// 💾 Insert session row
func (r *GoogleAuthRepository) InsertUserSession(ctx context.Context, arg sqlc.InsertUserSessionParams) (uuid.UUID, error) {
	return r.q.InsertUserSession(ctx, arg)
}

// 🔐 Insert refresh token row
func (r *GoogleAuthRepository) InsertRefreshToken(ctx context.Context, arg sqlc.InsertRefreshTokenParams) error {
	return r.q.InsertRefreshToken(ctx, arg)
}
