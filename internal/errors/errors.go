package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrProductNotFound is returned when a product is not found
	ErrKeyNotFound = errors.New("product key not found")
	// ErrProductDeactivated is returned when a product is deactivated
	ErrKeyDeactivated = errors.New("product key is deactivated")
	// ErrUnauthorized is returned when a user is not authorized to perform an action
	ErrUnauthorized = errors.New("you do not own this key")
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")

)
