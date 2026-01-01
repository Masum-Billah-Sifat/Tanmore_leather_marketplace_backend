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

// ✅ Parse a string to *uuid.UUID — return nil if blank or invalid
func ParsePtr(input string) *uuid.UUID {
	if input == "" {
		return nil
	}
	id, err := uuid.Parse(input)
	if err != nil {
		return nil
	}
	return &id
}
