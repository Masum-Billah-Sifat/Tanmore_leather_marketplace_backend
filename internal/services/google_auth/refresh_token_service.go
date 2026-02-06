// ------------------------------------------------------------
// 📁 File: internal/service/auth/refresh_token_service.go
// 🧠 This file contains the logic for secure refresh token rotation,
//     including session checks, user status checks, and token regeneration.

package google_auth

import (
	"context"
	"fmt"
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

type RefreshTokenOutput struct {
	AccessToken             string
	RefreshToken            string
	ExpiresIn               int
	IsSellerProfileApproved bool // ✅ add this
}

// 🔐 Service to manage token rotation
type RefreshTokenService struct {
	repo token_refresh.TokenRefreshRepoInterface
}

// 🚀 Constructor
func NewRefreshTokenService(repo token_refresh.TokenRefreshRepoInterface) *RefreshTokenService {
	return &RefreshTokenService{repo: repo}
}

func (s *RefreshTokenService) HandleRefreshTokenRotation(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	fmt.Println("🔁 [Refresh Flow] Got request")
	fmt.Println("🧾 Raw token:", input.RawToken)
	fmt.Println("🧾 User Agent:", input.UserAgent)
	fmt.Println("🧾 Device Fingerprint:", input.DeviceFingerprint)

	tokenHash := token.HashRefreshToken(input.RawToken)

	var (
		user    sqlc.User
		session sqlc.UserSession
		output  *RefreshTokenOutput
	)

	err := s.repo.WithTx(ctx, func(q *sqlc.Queries) error {
		refreshToken, err := q.GetRefreshTokenByHash(ctx, tokenHash)
		if err != nil || refreshToken.IsDeprecated || refreshToken.ExpiresAt.Before(timeutil.NowUTC()) {
			return errors.NewAuthError("invalid or expired refresh token")
		}

		user, err = q.GetUserByID(ctx, refreshToken.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived || user.IsBanned {
			return errors.NewAuthError("user is not allowed to refresh session")
		}

		session, err = q.GetSessionByIDAndUserID(ctx, sqlc.GetSessionByIDAndUserIDParams{
			ID:     refreshToken.SessionID,
			UserID: refreshToken.UserID,
		})
		if err != nil || session.IsRevoked || session.IsArchived {
			return errors.NewAuthError("session is invalid or revoked")
		}

		if session.UserAgent != input.UserAgent || session.DeviceFingerprint != input.DeviceFingerprint {
			return errors.NewAuthError("session fingerprint mismatch")
		}

		err = q.DeprecateRefreshTokenByID(ctx, sqlc.DeprecateRefreshTokenByIDParams{
			ID:               refreshToken.ID,
			IsDeprecated:     true,
			DeprecatedReason: sqlnull.String("rotated"),
			DeprecatedAt:     sqlnull.Time(timeutil.NowUTC()),
		})
		if err != nil {
			return errors.NewServerError("deprecating old token")
		}

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

		accessToken, err := token.GenerateAccessToken(user.ID, session.ID, user.CurrentMode, 1)
		if err != nil {
			return errors.NewServerError("generating access token")
		}

		// ✅ Include is_seller_profile_approved in output
		output = &RefreshTokenOutput{
			AccessToken:             accessToken,
			RefreshToken:            rawRefreshToken,
			ExpiresIn:               15 * 60,
			IsSellerProfileApproved: user.IsSellerProfileApproved, // ✅ NEW FIELD
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return output, nil
}
