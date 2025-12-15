// ------------------------------------------------------------
// 📁 File: internal/repository/user_mode_switch/user_mode_switch_repo.go
// 🧠 Defines interface for switching user mode and logging the switch

package user_mode_switch

import (
	"context"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserModeSwitchRepoInterface interface {
	WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error)
	UpdateUserCurrentMode(ctx context.Context, userID uuid.UUID, toMode string) error
	InsertUserModeSwitchLog(ctx context.Context, log sqlc.InsertUserModeSwitchLogParams) error
}
