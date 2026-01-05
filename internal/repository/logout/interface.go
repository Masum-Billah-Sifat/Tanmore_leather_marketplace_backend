// ------------------------------------------------------------
// 📁 File: internal/repository/logout/interface.go
// 🧠 This file defines the interface for LogoutRepoInterface
//     which provides all DB operations needed for logout flow.

package logout

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type LogoutRepoInterface interface {
	// 🔁 Transaction handler
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error

	// 🔎 Fetch session by ID and user ID
	GetSessionByIDAndUserID(ctx context.Context, sessionID, userID uuid.UUID) (sqlc.UserSession, error)

	// 🚫 Revoke the session
	RevokeUserSession(ctx context.Context, arg sqlc.RevokeUserSessionParams) error

	// 🗑️ Deprecate all refresh tokens for this session
	DeprecateRefreshTokensBySession(ctx context.Context, arg sqlc.DeprecateRefreshTokensBySessionParams) error
}
