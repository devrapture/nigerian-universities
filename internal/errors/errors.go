package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrKeyNotFound is returned when a product key is not found
	ErrKeyNotFound = errors.New("product key not found")
	// ErrKeyDeactivated is returned when a product key is deactivated
	ErrKeyDeactivated = errors.New("product key is deactivated")
	// ErrUnauthorized is returned when a user is not authorized to perform an action
	ErrUnauthorized = errors.New("you do not own this key")
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
)
