// user_service.go — Business logic for users

package service

import (
	"context"
	"database/sql"
	"tanmore_backend/internal/db/sqlc"
)

type UserService struct {
	queries *sqlc.Queries
}

func NewUserService(q *sqlc.Queries) *UserService {
	return &UserService{queries: q}
}

func (s *UserService) CreateUser(ctx context.Context, googleID, primaryEmail, displayName, profileImageURL string) (sqlc.User, error) {
	params := sqlc.CreateUserParams{
		GoogleID:     googleID,
		PrimaryEmail: primaryEmail,
		// DisplayName:     displayName,
		// ProfileImageUrl: profileImageURL,
		DisplayName: sql.NullString{
			String: displayName,
			Valid:  displayName != "",
		},
		ProfileImageUrl: sql.NullString{
			String: profileImageURL,
			Valid:  profileImageURL != "",
		},
	}
	return s.queries.CreateUser(ctx, params)
}
