// ------------------------------------------------------------
// 📁 File: internal/repository/token_refresh/interface.go
// 🧠 This file defines the interface for TokenRefreshRepoInterface
//     which provides all DB operations needed for refresh token rotation.

package token_refresh

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type TokenRefreshRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🔍 Lookup refresh token by hash
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (sqlc.UserRefreshToken, error)

	// 👤 Lookup user by ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)

	// 📲 Lookup session by session ID + user ID
	GetSessionByIDAndUserID(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (sqlc.UserSession, error)

	// 🚫 Deprecate old refresh token
	DeprecateRefreshTokenByID(ctx context.Context, arg sqlc.DeprecateRefreshTokenByIDParams) error

	// 🔁 Insert new refresh token
	InsertRefreshToken(ctx context.Context, arg sqlc.InsertRefreshTokenParams) error
}
