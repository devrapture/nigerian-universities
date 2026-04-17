package schema

import "time"

type KeyItem struct {
	ID         string     `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID     string     `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	KeyPrefix  string     `json:"key_prefix" example:"sk_live_"`
	IsActive   bool       `json:"is_active" example:"true"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type GeneratedKey struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Key       string    `json:"key" example:"sk_live_0123456789abcdef"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KeyCreateResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Store this key securely. It will not be shown again."`
	Data    GeneratedKey `json:"data"`
}

type KeyListResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Fetched all keys"`
	Data    []KeyItem      `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

type KeySuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Key deactivated successfully"`
}

type KeyBadRequestError struct {
	Code    string `json:"code" example:"BAD_REQUEST"`
	Message string `json:"message" example:"invalid key id"`
}

type KeyBadRequestResponse struct {
	Success bool                `json:"success" example:"false"`
	Error   *KeyBadRequestError `json:"error,omitempty"`
}

type KeyUnauthorizedError struct {
	Code    string `json:"code" example:"UNAUTHORIZED"`
	Message string `json:"message" example:"Invalid or expired token"`
}

type KeyUnauthorizedResponse struct {
	Success bool                 `json:"success" example:"false"`
	Error   *KeyUnauthorizedError `json:"error,omitempty"`
}

type KeyNotFoundError struct {
	Code    string `json:"code" example:"NOT_FOUND"`
	Message string `json:"message" example:"key not found"`
}

type KeyNotFoundResponse struct {
	Success bool             `json:"success" example:"false"`
	Error   *KeyNotFoundError `json:"error,omitempty"`
}

type KeyInternalServerError struct {
	Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Message string `json:"message" example:"failed to create api key"`
}

type KeyInternalServerErrorResponse struct {
	Success bool                       `json:"success" example:"false"`
	Error   *KeyInternalServerError `json:"error,omitempty"`
}
