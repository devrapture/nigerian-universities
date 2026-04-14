package handlers

import (
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
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.userService.GetGoogleAuthURL(c)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}
	user, jwtToken, err := h.userService.HandleGoogleCallback(c, code)
	if err != nil {
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

func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	var req LoginWithGoogleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}
	user, jwtToken, err := h.userService.HandleLoginWithGoogle(c.Request.Context(), req.ID, req.Email, req.Name, req.AvatarURL)
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

func (h *AuthHandler) GithubLogin(c *gin.Context) {
	url := h.userService.GetGithubAuthUrl(c)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GihubCallback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "code not found",
		})
	}

	user, jwtToken, err := h.userService.HandleGithubCallback(c.Request.Context(), code)
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
