package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/oauth2"
)

// OAuthService provides common OAuth functionality
type OAuthService struct {
	config *oauth2.Config
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(config *oauth2.Config) *OAuthService {
	return &OAuthService{
		config: config,
	}
}

// GetAuthURL returns the OAuth URL for authentication
func (s *OAuthService) GetAuthURL() string {
	state := generateRandomState()
	return s.config.AuthCodeURL(state)
}

// ExchangeCode exchanges an authorization code for a token
func (s *OAuthService) ExchangeCode(code string) (*oauth2.Token, error) {
	return s.config.Exchange(context.Background(), code)
}

// generateRandomState generates a random state for OAuth
func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
