// ------------------------------------------------------------
// 📁 File: internal/repository/google_auth/interface.go
// 🧠 This file defines the interface for GoogleAuthRepoInterface
//     which provides all DB operations needed for Google login flow.

package google_auth

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type GoogleAuthRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 👤 Fetch user by Google ID
	GetUserByGoogleID(ctx context.Context, googleID string) (sqlc.User, error)

	// ➕ Insert new user
	InsertUser(ctx context.Context, arg sqlc.InsertUserParams) (uuid.UUID, error)

	// 💾 Insert new session
	InsertUserSession(ctx context.Context, arg sqlc.InsertUserSessionParams) (uuid.UUID, error)

	// 🔐 Insert refresh token
	InsertRefreshToken(ctx context.Context, arg sqlc.InsertRefreshTokenParams) error
}
