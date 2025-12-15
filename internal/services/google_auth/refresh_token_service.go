// ------------------------------------------------------------
// 📁 File: internal/service/auth/refresh_token_service.go
// 🧠 This file contains the logic for secure refresh token rotation,
//     including session checks, user status checks, and token regeneration.

package google_auth

import (
	"context"
	"time"

	"tanmore_backend/internal/db/sqlc"
	"tanmore_backend/internal/repository/token_refresh"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 🔑 Input struct for refresh token handler
type RefreshTokenInput struct {
	RawToken          string
	UserAgent         string
	Platform          string
	DeviceFingerprint string
	IPAddress         string
}

// 📦 Output struct to return new tokens
type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // seconds
}

// 🔐 Service to manage token rotation
type RefreshTokenService struct {
	repo token_refresh.TokenRefreshRepoInterface
}

// 🚀 Constructor
func NewRefreshTokenService(repo token_refresh.TokenRefreshRepoInterface) *RefreshTokenService {
	return &RefreshTokenService{repo: repo}
}

// 🔁 Main refresh flow
func (s *RefreshTokenService) HandleRefreshTokenRotation(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	tokenHash := token.HashRefreshToken(input.RawToken)

	var (
		user    sqlc.User
		session sqlc.UserSession
		output  *RefreshTokenOutput
	)

	err := s.repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Find the refresh token
		refreshToken, err := q.GetRefreshTokenByHash(ctx, tokenHash)
		if err != nil {
			return errors.NewAuthError("invalid or expired refresh token")
		}
		if refreshToken.IsDeprecated || refreshToken.ExpiresAt.Before(timeutil.NowUTC()) {
			return errors.NewAuthError("refresh token is no longer valid")
		}

		// Step 2: Get user and validate
		user, err = q.GetUserByID(ctx, refreshToken.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived || user.IsBanned {
			return errors.NewAuthError("user is not allowed to refresh session")
		}

		// Step 3: Get session and validate
		session, err = q.GetSessionByIDAndUserID(ctx, sqlc.GetSessionByIDAndUserIDParams{
			ID:     refreshToken.SessionID,
			UserID: refreshToken.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("session")
		}
		if session.IsRevoked || session.IsArchived {
			return errors.NewAuthError("session is invalid or revoked")
		}

		// Step 4: Match headers
		if session.UserAgent != input.UserAgent || session.DeviceFingerprint != input.DeviceFingerprint {
			return errors.NewAuthError("session fingerprint mismatch")
		}

		// Step 5: Deprecate old token
		err = q.DeprecateRefreshTokenByID(ctx, sqlc.DeprecateRefreshTokenByIDParams{
			ID:               refreshToken.ID,
			IsDeprecated:     true,
			DeprecatedReason: sqlnull.String("rotated"),
			DeprecatedAt:     sqlnull.Time(timeutil.NowUTC()),
		})
		if err != nil {
			return errors.NewServerError("deprecating old token")
		}

		// Step 6: Insert new refresh token
		newRefreshID := uuid.New()
		rawRefreshToken, err := token.GenerateRefreshToken()
		if err != nil {
			return errors.NewServerError("generating refresh token")
		}

		err = q.InsertRefreshToken(ctx, sqlc.InsertRefreshTokenParams{
			ID:               newRefreshID,
			UserID:           user.ID,
			SessionID:        session.ID,
			TokenHash:        token.HashRefreshToken(rawRefreshToken),
			DeprecatedReason: sqlnull.String(""),
			IsDeprecated:     false,
			DeprecatedAt:     sqlnull.Time(time.Time{}),
			ExpiresAt:        timeutil.NowUTC().Add(90 * 24 * time.Hour),
			CreatedAt:        timeutil.NowUTC(),
		})
		if err != nil {
			return errors.NewServerError("inserting new refresh token")
		}

		// Step 7: Generate new access token
		accessToken, err := token.GenerateAccessToken(user.ID, session.ID, user.CurrentMode, 15)
		if err != nil {
			return errors.NewServerError("generating access token")
		}

		// ✅ Prepare final output
		output = &RefreshTokenOutput{
			AccessToken:  accessToken,
			RefreshToken: rawRefreshToken,
			ExpiresIn:    15 * 60, // 15 minutes
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return output, nil
}
