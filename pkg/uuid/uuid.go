package uuidutil

import "github.com/google/uuid"

// Generate a fresh UUID
func New() uuid.UUID {
	return uuid.New()
}

// Parse a UUID string safely
func Parse(input string) (uuid.UUID, error) {
	return uuid.Parse(input)
}

// Validate if string is a valid UUID
func IsValid(input string) bool {
	_, err := uuid.Parse(input)
	return err == nil
}
