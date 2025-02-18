package errors

import "errors"

var (
	// ErrKeyNotFound is returned when a key doesn't exist
	ErrKeyNotFound = errors.New("key not found")

	// ErrEmptyKey is returned when an empty key is provided
	ErrEmptyKey = errors.New("empty key not allowed")
)
