package types

import (
	"github.com/google/uuid"
)

// UUID represents a UUID type
type UUID uuid.UUID

// MustParseUUID parses a string into a UUID, panicking if invalid
func MustParseUUID(s string) UUID {
	return UUID(uuid.MustParse(s))
}

// String returns the string representation of the UUID
func (u UUID) String() string {
	return uuid.UUID(u).String()
} 