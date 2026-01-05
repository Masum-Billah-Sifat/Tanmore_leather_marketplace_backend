// ------------------------------------------------------------
// 📁 File: internal/service/logout/logout_service.go
// 🧠 Handles logout flow:
//     - Validates session existence
//     - Ensures it's not already revoked/archived
//     - Revokes session
//     - Deprecates refresh tokens
//     - Returns confirmation message

package logout

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/logout"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📦 Input from handler
type LogoutInput struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}

// ------------------------------------------------------------
// 📦 Output returned to handler
type LogoutResult struct {
	Message string
}

// ------------------------------------------------------------
// 📦 Dependencies
type LogoutServiceDeps struct {
	Repo repo.LogoutRepoInterface
}

// ------------------------------------------------------------
// 📦 Service definition
type LogoutService struct {
	Deps LogoutServiceDeps
}

// ------------------------------------------------------------
// 🛠️ Constructor
func NewLogoutService(deps LogoutServiceDeps) *LogoutService {
	return &LogoutService{Deps: deps}
}

// ------------------------------------------------------------
// 🚪 Logout handler
func (s *LogoutService) HandleLogout(ctx context.Context, input LogoutInput) (*LogoutResult, error) {
	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// 🔍 Step 1: Validate session
		session, err := q.GetSessionByIDAndUserID(ctx, sqlc.GetSessionByIDAndUserIDParams{
			ID:     input.SessionID,
			UserID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("session")
		}
		if session.IsRevoked || session.IsArchived {
			return errors.NewValidationError("session", "session already closed")
		}

		// 🔒 Step 2: Revoke session
		err = q.RevokeUserSession(ctx, sqlc.RevokeUserSessionParams{
			ID:        input.SessionID,
			UserID:    input.UserID,
			UpdatedAt: now,
			IsRevoked: true,
		})
		if err != nil {
			return errors.NewServerError("revoke session failed")
		}

		// 🗑️ Step 3: Deprecate all refresh tokens for this session
		err = q.DeprecateRefreshTokensBySession(ctx, sqlc.DeprecateRefreshTokensBySessionParams{
			UserID:           input.UserID,
			SessionID:        input.SessionID,
			DeprecatedAt:     sqlnull.Time(now),
			DeprecatedReason: sqlnull.String("manual_logout"),
			IsDeprecated:     true,
		})
		if err != nil {
			return errors.NewServerError("deprecate refresh tokens failed")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &LogoutResult{
		Message: "Logout successful",
	}, nil
}
