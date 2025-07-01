package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"jti-super-app-go/config"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type GoogleAuthService interface {
	GenerateAuthURL(state string) string
	GetUserInfo(code string) (*GoogleUserInfo, error)
}

type googleAuthService struct {
	oauthConfig *oauth2.Config
}

func NewGoogleAuthService(cfg *config.Config) GoogleAuthService {
	return &googleAuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.AppUrl + "/api/v1/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *googleAuthService) GenerateAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *googleAuthService) GetUserInfo(code string) (*GoogleUserInfo, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange code for token: " + err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.New("failed to get user info from Google: " + err.Error())
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		return nil, errors.New("failed to unmarshal user info: " + err.Error())
	}

	if !userInfo.VerifiedEmail {
		return nil, errors.New("email from Google is not verified")
	}

	return &userInfo, nil
}
