package handlers

import (
	"errors"
	"net/http"

	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

type GoogleLoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        UserPayload `json:"user"`
}

type GithubLoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        UserPayload `json:"user"`
}

type UserPayload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

type LoginWithGoogleRequest struct {
	IDToken string `json:"id_token"`
}

type LoginWithGithubRequest struct {
	AccessToken string `json:"access_token"`
}

// @Summary Google Login
// @Description Google Login
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} schema.AuthGoogleLoginResponse
// @Failure 400 {object} schema.AuthBadRequestResponse
// @Failure 500 {object} schema.AuthInternalServerErrorResponse
// @Router /auth/google [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url, err := h.userService.GetGoogleAuthURL(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to initialize google oauth")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// @Summary Google Callback
// @Description Google OAuth callback
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "OAuth code"
// @Param state query string true "OAuth state"
// @Success 200 {object} schema.AuthGoogleLoginResponse
// @Failure 400 {object} schema.AuthBadRequestResponse
// @Failure 500 {object} schema.AuthInternalServerErrorResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}
	if state == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "state not found")
		return
	}
	user, jwtToken, err := h.userService.HandleGoogleCallback(c, code, state)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOAuthState) {
			utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "invalid oauth state")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Authentication failed")
		return
	}

	response := GoogleLoginResponse{
		AccessToken: jwtToken,
		User:        toUserPayload(user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", response, nil)
}

func toUserPayload(u *model.User) UserPayload {
	if u == nil {
		return UserPayload{}
	}
	return UserPayload{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}
}

// @Summary Login with Google token
// @Description Login with Google ID token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body schema.LoginWithGoogleRequest true "Google login payload"
// @Success 200 {object} schema.AuthGoogleLoginResponse
// @Failure 400 {object} schema.AuthBadRequestResponse
// @Failure 500 {object} schema.AuthInternalServerErrorResponse
// @Router /auth/google/login [post]
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	var req LoginWithGoogleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}
	if req.IDToken == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "id_token is required")
		return
	}
	user, jwtToken, err := h.userService.HandleLoginWithGoogle(c.Request.Context(), req.IDToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response := GoogleLoginResponse{
		AccessToken: jwtToken,
		User:        toUserPayload(user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", response, nil)
}

// @Summary Github Login
// @Description Github Login
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} schema.AuthGithubLoginResponse
// @Failure 400 {object} schema.AuthBadRequestResponse
// @Failure 500 {object} schema.AuthInternalServerErrorResponse
// @Router /auth/github [get]
func (h *AuthHandler) GithubLogin(c *gin.Context) {
	url, err := h.userService.GetGithubAuthUrl(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to initialize github oauth")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// @Summary Github Callback
// @Description Github OAuth callback
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "OAuth code"
// @Param state query string true "OAuth state"
// @Success 200 {object} schema.AuthGithubLoginResponse
// @Failure 400 {object} schema.AuthBadRequestResponse
// @Failure 500 {object} schema.AuthInternalServerErrorResponse
// @Router /auth/github/callback [get]
func (h *AuthHandler) GithubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "code not found",
		})
		return
	}
	if state == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "state not found")
		return
	}

	user, jwtToken, err := h.userService.HandleGithubCallback(c.Request.Context(), code, state)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOAuthState) {
			utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "invalid oauth state")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response := GithubLoginResponse{
		AccessToken: jwtToken,
		User:        toUserPayload(user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", response, nil)
}

// @Summary Login with Github token
// @Description Login with Github access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body schema.LoginWithGithubRequest true "Github login payload"
// @Success 200 {object} schema.AuthGithubLoginResponse
// @Failure 400 {object} schema.AuthBadRequestResponse
// @Failure 500 {object} schema.AuthInternalServerErrorResponse
// @Router /auth/github/login [post]
func (h *AuthHandler) LoginWithGithub(c *gin.Context) {
	var req LoginWithGithubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}
	if req.AccessToken == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "access_token is required")
		return
	}

	user, jwtToken, err := h.userService.HandleLoginWithGithub(c.Request.Context(), req.AccessToken)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response := GithubLoginResponse{
		AccessToken: jwtToken,
		User:        toUserPayload(user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Success", response, nil)
}
