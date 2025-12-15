package errors

import "fmt"

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
