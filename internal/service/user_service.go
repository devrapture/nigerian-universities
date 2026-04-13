package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type UserService interface {
	GetGoogleAuthURL(ctx context.Context) string
	HandleGoogleCallback(ctx context.Context, code string) (*model.User, string, error)
	HandleLoginWithGoogle(ctx context.Context, id, email, name, picture string) (*model.User, string, error)
}

type userService struct {
	repo              repositories.UserRepository
	googleOAuthConfig *oauth2.Config
	githubOAuthConfig *oauth2.Config
	cfg               *config.Config
}

func NewUserService(cfg *config.Config, repo repositories.UserRepository) UserService {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.GOOGLE_CLIENT_ID,
		ClientSecret: cfg.GOOGLE_CLIENT_SECRET,
		RedirectURL:  cfg.GOOGLE_REDIRECT_URL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	githubConfig := &oauth2.Config{
		ClientID:     cfg.GITHUB_CLIENT_ID,
		ClientSecret: cfg.GITHUB_CLIENT_SECRET,
		RedirectURL:  cfg.GITHUB_REDIRECT_URL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
	return &userService{
		repo:              repo,
		googleOAuthConfig: googleConfig,
		githubOAuthConfig: githubConfig,
		cfg:               cfg,
	}
}

func (s *userService) HandleGoogleCallback(ctx context.Context, code string) (*model.User, string, error) {
	token, err := s.googleOAuthConfig.Exchange(ctx, code)

	if err != nil {
		return nil, "", errors.New("failed to exchange token")
	}
	client := s.googleOAuthConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("failed to get user info")
	}

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, "", errors.New("failed to parse user info")
	}
	user, err := s.repo.FindOrCreateUser(ctx, userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}
	jwtToken, err := utils.GenerateJwt(user.ID, user.Email, s.cfg)

	if err != nil {
		return nil, "", errors.New("failed to generate jwt")
	}
	return user, jwtToken, nil

}

func (s *userService) GetGoogleAuthURL(ctx context.Context) string {
	url := s.googleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

func (s *userService) HandleLoginWithGoogle(ctx context.Context, id, email, name, picture string) (*model.User, string, error) {
	user, err := s.repo.FindOrCreateUser(ctx, id, email, name, picture)
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}
	jwtToken, err := utils.GenerateJwt(user.ID, user.Email, s.cfg)

	if err != nil {
		return nil, "", errors.New("failed to generate jwt")
	}
	return user, jwtToken, nil

}
