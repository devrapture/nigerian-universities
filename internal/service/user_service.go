package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var ErrInvalidOAuthState = errors.New("invalid oauth state")

type UserService interface {
	GetGoogleAuthURL(ctx context.Context) (string, error)
	HandleGoogleCallback(ctx context.Context, code, state string) (*model.User, string, error)
	HandleLoginWithGoogle(ctx context.Context, id, email, name, picture string) (*model.User, string, error)
	HandleLoginWithGithub(ctx context.Context, id, email, name, picture string) (*model.User, string, error)
	GetGithubAuthUrl(ctx context.Context) (string, error)
	HandleGithubCallback(ctx context.Context, code, state string) (*model.User, string, error)
}

type userService struct {
	repo              repositories.UserRepository
	googleOAuthConfig *oauth2.Config
	githubOAuthConfig *oauth2.Config
	cfg               *config.Config
	stateStore        *oauthStateStore
}

type oauthStateStore struct {
	mu     sync.Mutex
	tokens map[string]oauthStateEntry
	ttl    time.Duration
}

type oauthStateEntry struct {
	provider  string
	expiresAt time.Time
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
		stateStore: &oauthStateStore{
			tokens: make(map[string]oauthStateEntry),
			ttl:    10 * time.Minute,
		},
	}
}

func (s *userService) HandleGoogleCallback(ctx context.Context, code, state string) (*model.User, string, error) {
	if err := s.stateStore.ValidateAndConsume("google", state); err != nil {
		return nil, "", err
	}

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
	user, err := s.repo.FindOrCreateUser(ctx, userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture, "google")
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}
	jwtToken, err := utils.GenerateJwt(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", errors.New("failed to generate jwt")
	}
	return user, jwtToken, nil
}

func (s *userService) GetGoogleAuthURL(ctx context.Context) (string, error) {
	state, err := s.stateStore.Generate("google")
	if err != nil {
		return "", err
	}

	return s.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *userService) GetGithubAuthUrl(ctx context.Context) (string, error) {
	state, err := s.stateStore.Generate("github")
	if err != nil {
		return "", err
	}

	return s.githubOAuthConfig.AuthCodeURL(state), nil
}

func (s *userService) HandleLoginWithGoogle(ctx context.Context, id, email, name, picture string) (*model.User, string, error) {
	user, err := s.repo.FindOrCreateUser(ctx, id, email, name, picture, "google")
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}
	jwtToken, err := utils.GenerateJwt(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", errors.New("failed to generate jwt")
	}
	return user, jwtToken, nil
}

func (s *userService) HandleGithubCallback(ctx context.Context, code, state string) (*model.User, string, error) {
	if err := s.stateStore.ValidateAndConsume("github", state); err != nil {
		return nil, "", err
	}

	token, err := s.githubOAuthConfig.Exchange(ctx, code)
	if err != nil {
		log.Println("failed to exchange token", err)
		return nil, "", errors.New("failed to exchange token")
	}

	client := s.githubOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, "", errors.New("failed to get user info")
	}

	defer resp.Body.Close()
	var userInfo struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, "", errors.New("failed to parse user info")
	}

	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()

			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}

			if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary {
						userInfo.Email = e.Email
						break
					}
				}
			}
		}
	}
	// find or create user
	user, err := s.repo.FindOrCreateUser(ctx, strconv.Itoa(userInfo.ID), userInfo.Email, userInfo.Name, userInfo.AvatarURL, "github")
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}

	jwtToken, err := utils.GenerateJwt(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", errors.New("failed to generate jwt")
	}
	return user, jwtToken, nil
}

func (s *oauthStateStore) Generate(provider string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("failed to generate oauth state")
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked(time.Now())
	s.tokens[token] = oauthStateEntry{
		provider:  provider,
		expiresAt: time.Now().Add(s.ttl),
	}

	return token, nil
}

func (s *oauthStateStore) ValidateAndConsume(provider, token string) error {
	if token == "" {
		return ErrInvalidOAuthState
	}

	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked(now)

	entry, ok := s.tokens[token]
	if !ok || entry.provider != provider || now.After(entry.expiresAt) {
		delete(s.tokens, token)
		return ErrInvalidOAuthState
	}

	delete(s.tokens, token)
	return nil
}

func (s *oauthStateStore) cleanupExpiredLocked(now time.Time) {
	for token, entry := range s.tokens {
		if now.After(entry.expiresAt) {
			delete(s.tokens, token)
		}
	}
}

func (s *userService) HandleLoginWithGithub(ctx context.Context, id, email, name, picture string) (*model.User, string, error) {
	user, err := s.repo.FindOrCreateUser(ctx, id, email, name, picture, "github")
	if err != nil {
		return nil, "", errors.New("failed to create user")
	}
	jwtToken, err := utils.GenerateJwt(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", errors.New("failed to generate jwt")
	}
	return user, jwtToken, nil
}
