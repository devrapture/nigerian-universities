package schema

type AuthUser struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Name      string `json:"name" example:"John Doe"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.png"`
}

type AuthTokenData struct {
	AccessToken string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User        AuthUser `json:"user"`
}

type AuthGoogleLoginResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Success"`
	Data    AuthTokenData `json:"data"`
}

type AuthGithubLoginResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Success"`
	Data    AuthTokenData `json:"data"`
}

type LoginWithGoogleRequest struct {
	IDToken string `json:"id_token" example:"google-id-token"`
}

type LoginWithGithubRequest struct {
	AccessToken string `json:"access_token" example:"github-access-token"`
}

type AuthBadRequestError struct {
	Code    string `json:"code" example:"BAD_REQUEST"`
	Message string `json:"message" example:"invalid request body"`
}

type AuthBadRequestResponse struct {
	Success bool                `json:"success" example:"false"`
	Error   *AuthBadRequestError `json:"error,omitempty"`
}

type AuthInternalServerError struct {
	Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Message string `json:"message" example:"Authentication failed"`
}

type AuthInternalServerErrorResponse struct {
	Success bool                        `json:"success" example:"false"`
	Error   *AuthInternalServerError `json:"error,omitempty"`
}
