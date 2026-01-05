// ------------------------------------------------------------
// 📁 File: internal/service/google_auth/google_auth_service.go
// 🧠 Handles Google Login / Registration workflow.
//     - Verifies Google ID token
//     - Creates user if new
//     - Creates user session
//     - Creates refresh token
//     - Issues access token
//     - Returns unified response object

package google_auth

import (
	"context"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/google_auth"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	tokenutil "tanmore_backend/pkg/token"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📦 Input from handler
type GoogleLoginInput struct {
	IDToken           string
	UserAgent         string
	Platform          string
	DeviceFingerprint string
	IPAddress         string
}

// ------------------------------------------------------------
// 📦 Result returned to handler
type GoogleLoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         GoogleLoginResultUser
}

// this is what i got rid of
//   struct {
// 		ID    uuid.UUID
// 		Name  string
// 		Email string
// 		Image string
// 	}

type GoogleLoginResultUser struct {
	ID                      uuid.UUID
	Name                    string
	Email                   string
	Image                   string
	IsSellerProfileApproved bool // ✅ Add this
}

// ------------------------------------------------------------
// 📦 Dependencies
type GoogleAuthServiceDeps struct {
	Repo repo.GoogleAuthRepoInterface
}

// ------------------------------------------------------------
// 📦 Service definition
type GoogleAuthService struct {
	Deps GoogleAuthServiceDeps
}

// ------------------------------------------------------------
// 🛠️ Constructor
func NewGoogleAuthService(deps GoogleAuthServiceDeps) *GoogleAuthService {
	return &GoogleAuthService{Deps: deps}
}

func (s *GoogleAuthService) Start(ctx context.Context, input GoogleLoginInput) (*GoogleLoginResult, error) {
	now := timeutil.NowUTC()

	googlePayload, err := tokenutil.VerifyGoogleIDToken(input.IDToken)
	if err != nil {
		return nil, errors.NewValidationError("id_token", "invalid google token")
	}

	googleID := googlePayload.Sub
	email := googlePayload.Email
	name := googlePayload.Name
	image := googlePayload.Picture

	var userID uuid.UUID
	var sessionID uuid.UUID
	var rawRefreshToken string
	var isSellerProfileApproved bool // ✅ track this

	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		existingUser, err := q.GetUserByGoogleID(ctx, googleID)
		if err == nil {
			// Existing user
			if existingUser.IsArchived || existingUser.IsBanned {
				return errors.NewAuthError("account is not allowed to login")
			}
			userID = existingUser.ID
			isSellerProfileApproved = existingUser.IsSellerProfileApproved // ✅ set from DB
		} else {
			// New user
			newUserID := uuidutil.New()
			userID = newUserID
			isSellerProfileApproved = false // ✅ default

			_, err := q.InsertUser(ctx, sqlc.InsertUserParams{
				ID:                      newUserID,
				GoogleID:                googleID,
				PrimaryEmail:            email,
				DisplayName:             sqlnull.String("some_dummy_stuff"),
				ProfileImageUrl:         sqlnull.String("some bunch of stuff_for_now"),
				IsArchived:              false,
				IsBanned:                false,
				IsMuted:                 false,
				CurrentMode:             "customer",
				IsSellerProfileApproved: false,
				IsSellerProfileCreated:  false,
				CreatedAt:               now,
				UpdatedAt:               now,
			})
			if err != nil {
				return errors.NewTableError("users.insert", err.Error())
			}
		}

		sessionID = uuidutil.New()
		_, err = q.InsertUserSession(ctx, sqlc.InsertUserSessionParams{
			ID:                sessionID,
			UserID:            userID,
			IpAddress:         input.IPAddress,
			UserAgent:         input.UserAgent,
			DeviceFingerprint: input.DeviceFingerprint,
			IsRevoked:         false,
			IsArchived:        false,
			CreatedAt:         now,
			UpdatedAt:         now,
		})
		if err != nil {
			return errors.NewTableError("user_sessions.insert", err.Error())
		}

		refreshID := uuidutil.New()
		rawRefreshToken, err = tokenutil.GenerateRefreshToken()
		if err != nil {
			return errors.NewServerError("generate refresh token")
		}
		refreshHash := tokenutil.HashRefreshToken(rawRefreshToken)

		err = q.InsertRefreshToken(ctx, sqlc.InsertRefreshTokenParams{
			ID:               refreshID,
			UserID:           userID,
			SessionID:        sessionID,
			TokenHash:        refreshHash,
			DeprecatedReason: sqlnull.String(""),
			IsDeprecated:     false,
			DeprecatedAt:     sqlnull.Time(time.Time{}),
			ExpiresAt:        now.Add(90 * 24 * time.Hour),
			CreatedAt:        now,
		})
		if err != nil {
			return errors.NewTableError("user_refresh_tokens.insert", err.Error())
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := tokenutil.GenerateAccessToken(userID, sessionID, "customer", 1)
	if err != nil {
		return nil, errors.NewServerError("generate access token")
	}

	return &GoogleLoginResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresIn:    900,
		User: GoogleLoginResultUser{
			ID:                      userID,
			Name:                    name,
			Email:                   email,
			Image:                   image,
			IsSellerProfileApproved: isSellerProfileApproved, // ✅ include it
		},
	}, nil
}

// // ------------------------------------------------------------
// // 🚀 Main Logic
// func (s *GoogleAuthService) Start(ctx context.Context, input GoogleLoginInput) (*GoogleLoginResult, error) {
// 	now := timeutil.NowUTC()

// 	// ------------------------------------------------------------
// 	// Step 1: Verify ID Token with Google
// 	googlePayload, err := tokenutil.VerifyGoogleIDToken(input.IDToken)
// 	if err != nil {
// 		return nil, errors.NewValidationError("id_token", "invalid google token")
// 	}

// 	googleID := googlePayload.Sub
// 	email := googlePayload.Email
// 	name := googlePayload.Name
// 	image := googlePayload.Picture

// 	var userID uuid.UUID
// 	var sessionID uuid.UUID
// 	var rawRefreshToken string

// 	// ------------------------------------------------------------
// 	// Step 2: Everything inside a transaction
// 	err = s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {

// 		// ------------------------------------------------------------
// 		// Step 2.1: Check if user exists
// 		existingUser, err := q.GetUserByGoogleID(ctx, googleID)
// 		if err == nil {
// 			// Found user
// 			if existingUser.IsArchived || existingUser.IsBanned {
// 				return errors.NewAuthError("account is not allowed to login")
// 			}
// 			userID = existingUser.ID
// 		} else {
// 			// ------------------------------------------------------------
// 			// Step 2.2: Insert new user
// 			newUserID := uuidutil.New()
// 			userID = newUserID

// 			_, err := q.InsertUser(ctx, sqlc.InsertUserParams{
// 				ID:                      newUserID,
// 				GoogleID:                googleID,
// 				PrimaryEmail:            email,
// 				DisplayName:             sqlnull.String("some_dummy_stuff"),
// 				ProfileImageUrl:         sqlnull.String("some bunch of stuff_for_now"),
// 				IsArchived:              false,
// 				IsBanned:                false,
// 				IsMuted:                 false,
// 				CurrentMode:             "customer",
// 				IsSellerProfileApproved: false,
// 				IsSellerProfileCreated:  false,
// 				CreatedAt:               now,
// 				UpdatedAt:               now,
// 			})
// 			if err != nil {
// 				return errors.NewTableError("users.insert", err.Error())
// 			}
// 		}

// 		// ------------------------------------------------------------
// 		// Step 2.3: Create new session
// 		sessionID = uuidutil.New()
// 		_, err = q.InsertUserSession(ctx, sqlc.InsertUserSessionParams{
// 			ID:                sessionID,
// 			UserID:            userID,
// 			IpAddress:         input.IPAddress,
// 			UserAgent:         input.UserAgent,
// 			DeviceFingerprint: input.DeviceFingerprint,
// 			IsRevoked:         false,
// 			IsArchived:        false,
// 			CreatedAt:         now,
// 			UpdatedAt:         now,
// 		})
// 		if err != nil {
// 			return errors.NewTableError("user_sessions.insert", err.Error())
// 		}

// 		// ------------------------------------------------------------
// 		// Step 2.4: Create refresh token
// 		refreshID := uuidutil.New()
// 		rawRefreshToken, err = tokenutil.GenerateRefreshToken()
// 		if err != nil {
// 			return errors.NewServerError("generate refresh token")
// 		}

// 		refreshHash := tokenutil.HashRefreshToken(rawRefreshToken)

// 		err = q.InsertRefreshToken(ctx, sqlc.InsertRefreshTokenParams{
// 			ID:               refreshID,
// 			UserID:           userID,
// 			SessionID:        sessionID,
// 			TokenHash:        refreshHash,
// 			DeprecatedReason: sqlnull.String(""),
// 			IsDeprecated:     false,
// 			DeprecatedAt:     sqlnull.Time(time.Time{}),
// 			ExpiresAt:        now.Add(90 * 24 * time.Hour),
// 			CreatedAt:        now,
// 		})
// 		if err != nil {
// 			return errors.NewTableError("user_refresh_tokens.insert", err.Error())
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	// ------------------------------------------------------------
// 	// Step 3: Generate Access Token
// 	accessToken, err := tokenutil.GenerateAccessToken(userID, sessionID, "customer", 1)
// 	if err != nil {
// 		return nil, errors.NewServerError("generate access token")
// 	}

// 	// ------------------------------------------------------------
// 	// Step 4: Prepare response
// 	result := &GoogleLoginResult{
// 		AccessToken:  accessToken,
// 		RefreshToken: rawRefreshToken,
// 		ExpiresIn:    900,
// 	}
// 	result.User.ID = userID
// 	result.User.Name = name
// 	result.User.Email = email
// 	result.User.Image = image

// 	return result, nil
// }
