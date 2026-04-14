package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type GoogleIdentityClaims struct {
	Subject string
	Email   string
	Name    string
	Picture string
}

type GithubIdentityClaims struct {
	ID      string
	Email   string
	Name    string
	Picture string
}

func VerifyGoogleIDToken(ctx context.Context, token, audience string) (*GoogleIdentityClaims, error) {
	if token == "" {
		return nil, errors.New("missing google id token")
	}

	payload, err := idtoken.Validate(ctx, token, audience)
	if err != nil {
		return nil, err
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if payload.Subject == "" || email == "" {
		return nil, errors.New("google id token missing required claims")
	}

	return &GoogleIdentityClaims{
		Subject: payload.Subject,
		Email:   email,
		Name:    name,
		Picture: picture,
	}, nil
}

func VerifyGithubAccessToken(ctx context.Context, accessToken string) (*GithubIdentityClaims, error) {
	if accessToken == "" {
		return nil, errors.New("missing github access token")
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch github user")
	}

	var userInfo struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Login     string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()

			if emailResp.StatusCode == http.StatusOK {
				var emails []struct {
					Email      string `json:"email"`
					Primary    bool   `json:"primary"`
					Verified   bool   `json:"verified"`
					Visibility string `json:"visibility"`
				}

				if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
					for _, e := range emails {
						if e.Primary && e.Verified {
							userInfo.Email = e.Email
							break
						}
					}
				}
			} else {
				return nil, errors.New("failed to fetch github user emails")
			}
		}
	}

	if userInfo.Name == "" {
		userInfo.Name = userInfo.Login
	}

	if userInfo.ID == 0 || userInfo.Email == "" {
		return nil, errors.New("github access token missing required claims")
	}

	return &GithubIdentityClaims{
		ID:      strconv.Itoa(userInfo.ID),
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.AvatarURL,
	}, nil
}
