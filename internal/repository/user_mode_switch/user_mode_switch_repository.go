// ------------------------------------------------------------
// 📁 File: internal/repository/user_mode_switch/user_mode_switch_repo_impl.go
// 🧠 Implements database logic for user mode switching using sqlc queries

package user_mode_switch

import (
	"context"
	"database/sql"
	"time"

	"tanmore_backend/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserModeSwitchRepository struct {
	db *sql.DB
}

func NewUserModeSwitchRepository(db *sql.DB) *UserModeSwitchRepository {
	return &UserModeSwitchRepository{db: db}
}

// ✅ Transaction wrapper
func (r *UserModeSwitchRepository) WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := sqlc.New(tx)

	if err := fn(q); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// 🔍 Get full user by ID
func (r *UserModeSwitchRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (sqlc.User, error) {
	q := sqlc.New(r.db)

	return q.GetUserByID(ctx, userID)
}

// ✏️ Update user mode (customer <-> seller)
func (r *UserModeSwitchRepository) UpdateUserCurrentMode(ctx context.Context, userID uuid.UUID, toMode string) error {
	q := sqlc.New(r.db)

	// Fix 2: Use time.Time directly (not sql.NullTime)
	return q.UpdateUserCurrentMode(ctx, sqlc.UpdateUserCurrentModeParams{
		ID:          userID,
		CurrentMode: toMode,
		UpdatedAt:   time.Now().UTC(),
	})
}

// 🧾 Insert mode switch log row
func (r *UserModeSwitchRepository) InsertUserModeSwitchLog(ctx context.Context, log sqlc.InsertUserModeSwitchLogParams) error {
	q := sqlc.New(r.db)
	return q.InsertUserModeSwitchLog(ctx, log)
}
