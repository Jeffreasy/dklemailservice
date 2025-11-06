package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UUIDArray is a custom type for handling PostgreSQL UUID arrays
type UUIDArray []uuid.UUID

// Scan implements sql.Scanner interface for reading UUID arrays from database
func (a *UUIDArray) Scan(value interface{}) error {
	if value == nil {
		*a = UUIDArray{}
		return nil
	}

	// PostgreSQL returns UUID arrays as pq.GenericArray
	switch v := value.(type) {
	case []byte:
		// Handle byte array format
		return a.scanFromBytes(v)
	case string:
		// Handle string format
		return a.scanFromString(v)
	default:
		// Try to use pq.Array to handle the native PostgreSQL array
		var stringArray pq.StringArray
		if err := stringArray.Scan(value); err != nil {
			return fmt.Errorf("failed to scan UUID array: %v", err)
		}
		return a.scanFromStringArray(stringArray)
	}
}

// scanFromBytes parses UUID array from byte slice
func (a *UUIDArray) scanFromBytes(b []byte) error {
	return a.scanFromString(string(b))
}

// scanFromString parses UUID array from string representation
func (a *UUIDArray) scanFromString(s string) error {
	// PostgreSQL array format: {uuid1,uuid2,uuid3}
	s = strings.TrimSpace(s)
	if s == "" || s == "{}" {
		*a = UUIDArray{}
		return nil
	}

	// Remove braces
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = UUIDArray{}
		return nil
	}

	// Split by comma
	parts := strings.Split(s, ",")
	result := make([]uuid.UUID, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		uid, err := uuid.Parse(part)
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %s - %v", part, err)
		}
		result = append(result, uid)
	}

	*a = result
	return nil
}

// scanFromStringArray converts pq.StringArray to UUIDArray
func (a *UUIDArray) scanFromStringArray(sa pq.StringArray) error {
	result := make([]uuid.UUID, 0, len(sa))
	for _, s := range sa {
		uid, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %s - %v", s, err)
		}
		result = append(result, uid)
	}
	*a = result
	return nil
}

// Value implements driver.Valuer interface for writing UUID arrays to database
func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Convert to string array format for PostgreSQL
	strs := make([]string, len(a))
	for i, uid := range a {
		strs[i] = uid.String()
	}

	// Use pq.Array for proper PostgreSQL array encoding
	return pq.Array(strs).Value()
}

// MarshalJSON implements json.Marshaler for JSON serialization
func (a UUIDArray) MarshalJSON() ([]byte, error) {
	if len(a) == 0 {
		return []byte("[]"), nil
	}

	strs := make([]string, len(a))
	for i, uid := range a {
		strs[i] = fmt.Sprintf(`"%s"`, uid.String())
	}
	return []byte(fmt.Sprintf("[%s]", strings.Join(strs, ","))), nil
}

// UnmarshalJSON implements json.Unmarshaler for JSON deserialization
func (a *UUIDArray) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		*a = UUIDArray{}
		return nil
	}

	// Handle empty array
	s := strings.TrimSpace(string(data))
	if s == "[]" {
		*a = UUIDArray{}
		return nil
	}

	// Parse JSON array
	s = strings.Trim(s, "[]")
	if s == "" {
		*a = UUIDArray{}
		return nil
	}

	// Split by comma and parse each UUID
	parts := strings.Split(s, ",")
	result := make([]uuid.UUID, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, `"`)
		if part == "" {
			continue
		}

		uid, err := uuid.Parse(part)
		if err != nil {
			return fmt.Errorf("invalid UUID in JSON array: %s - %v", part, err)
		}
		result = append(result, uid)
	}

	*a = result
	return nil
}
