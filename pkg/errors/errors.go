package errors

import (
	"fmt"
)

// 🔎 Validation error (missing fields, invalid types, etc.)
func NewValidationError(field string, message string) error {
	return fmt.Errorf("validation error: field '%s' %s", field, message)
}

// 🔐 Authentication / authorization failure
func NewAuthError(message string) error {
	return fmt.Errorf("auth error: %s", message)
}

// 🚫 Rate limit exceeded
func NewRateLimitError(action string) error {
	return fmt.Errorf("rate limit reached for action: %s", action)
}

// 🔍 Not found
func NewNotFoundError(entity string) error {
	return fmt.Errorf("not found: %s", entity)
}

// ⚠️ Conflict error (duplicate, already exists, etc.)
func NewConflictError(message string) error {
	return fmt.Errorf("conflict: %s", message)
}

// 💥 Internal server error with context
func NewServerError(context string) error {
	return fmt.Errorf("server error: %s", context)
}

// 🗂️ Database table-specific errors
func NewTableError(table string, reason string) error {
	return fmt.Errorf("table '%s' error: %s", table, reason)
}

// IsCustomCode checks for a specific message used as a stand-in error code.
func IsCustomCode(err error, code string) bool {
	return err != nil && err.Error() == code
}

// Auth error codes/messages used across backend
const (
	ErrMissingAccessToken  = "missing access token"
	ErrInvalidAccessToken  = "invalid or expired access token"
	ErrInvalidUserIDFormat = "invalid user ID format"
	ErrUserNotFound        = "user not found"
	ErrUserArchived        = "user is archived"
	ErrUserBanned          = "user is banned"
	ErrForbiddenCustomer   = "only accessible to customers"
	ErrForbiddenSeller     = "only accessible to sellers"
)

// ✅ New style explicit constructor wrappers
func ErrAuthMissingToken() error {
	return NewAuthError(ErrMissingAccessToken)
}
func ErrAuthInvalidToken() error {
	return NewAuthError(ErrInvalidAccessToken)
}
func ErrAuthInvalidUserID() error {
	return NewAuthError(ErrInvalidUserIDFormat)
}
func ErrAuthUserNotFound() error {
	return NewAuthError(ErrUserNotFound)
}
func ErrAuthArchivedUser() error {
	return NewAuthError(ErrUserArchived)
}
func ErrAuthBannedUser() error {
	return NewAuthError(ErrUserBanned)
}
func ErrAuthOnlyCustomer() error {
	return NewAuthError(ErrForbiddenCustomer)
}
func ErrAuthOnlySeller() error {
	return NewAuthError(ErrForbiddenSeller)
}
