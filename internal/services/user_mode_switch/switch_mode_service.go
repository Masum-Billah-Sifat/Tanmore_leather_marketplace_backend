// ------------------------------------------------------------
// 📁 File: internal/services/user_mode_switch/switch_mode_service.go
// 🧠 This file implements the service logic for switching a user's current mode
//     (customer ↔ seller) with validations, DB updates, history logging, and token generation.

package user_mode_switch

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/user_mode_switch"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/timeutil"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📥 Input struct for switching mode
type SwitchModeInput struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	FromMode  string
	ToMode    string
}

// 📤 Output struct for response
type SwitchModeOutput struct {
	AccessToken string
	ExpiresIn   int // in seconds
	Mode        string
}

// 🧠 Service struct
type SwitchModeService struct {
	repo repo.UserModeSwitchRepoInterface
}

// 🛠️ Constructor
func NewSwitchModeService(repo repo.UserModeSwitchRepoInterface) *SwitchModeService {
	return &SwitchModeService{repo: repo}
}

// 🚀 Main flow to handle switching
func (s *SwitchModeService) Handle(ctx context.Context, input SwitchModeInput) (*SwitchModeOutput, error) {
	var user sqlc.User

	err := s.repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Reject if same mode
		if input.FromMode == input.ToMode {
			return errors.NewValidationError("to_mode", "You are already in this mode")
		}

		// Step 2: Fetch user
		u, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if u.IsArchived {
			return errors.NewAuthError("user is archived")
		}
		if u.IsBanned {
			return errors.NewAuthError("user is banned")
		}

		user = u // store for access token later

		// Step 3: If switching to seller, validate seller preconditions
		if input.ToMode == "seller" {
			if !u.IsSellerProfileCreated {
				return errors.NewValidationError("to_mode", "Seller profile not created")
			}
			if !u.IsSellerProfileApproved {
				return errors.NewValidationError("to_mode", "Seller profile not approved by admin")
			}
		}

		// Step 4: Update users table
		err = q.UpdateUserCurrentMode(ctx, sqlc.UpdateUserCurrentModeParams{
			ID:          input.UserID,
			CurrentMode: input.ToMode,
			UpdatedAt:   timeutil.NowUTC(),
		})
		if err != nil {
			return errors.NewServerError("updating user mode")
		}

		// Step 5: Insert into user_mode_history
		err = q.InsertUserModeSwitchLog(ctx, sqlc.InsertUserModeSwitchLogParams{
			ID:         uuid.New(),
			UserID:     input.UserID,
			FromMode:   input.FromMode,
			ToMode:     input.ToMode,
			SwitchedAt: timeutil.NowUTC(),
		})
		if err != nil {
			return errors.NewServerError("logging mode switch")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Step 6: Generate new access token
	accessToken, err := token.GenerateAccessToken(user.ID, input.SessionID, input.ToMode, 15)
	if err != nil {
		return nil, errors.NewServerError("generating access token")
	}

	// ✅ Final output
	return &SwitchModeOutput{
		AccessToken: accessToken,
		ExpiresIn:   15 * 60,
		Mode:        input.ToMode,
	}, nil
}
