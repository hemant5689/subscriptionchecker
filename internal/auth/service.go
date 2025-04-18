package auth

import (
	"context"

	"golang.org/x/oauth2"
)

// Service provides OAuth2 authentication functionality
type Service struct {
	config *oauth2.Config
}

// NewService creates a new auth service with the given OAuth2 config
func NewService(cfg *oauth2.Config) *Service {
	return &Service{config: cfg}
}

// GetAuthURL returns the authorization URL for the OAuth2 flow
func (s *Service) GetAuthURL() string {
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("redirect_uri", s.config.RedirectURL),
	}
	return s.config.AuthCodeURL("state-token", opts...)
}

// ExchangeCode exchanges an authorization code for an OAuth2 token
func (s *Service) ExchangeCode(code string) (*oauth2.Token, error) {
	return s.config.Exchange(context.Background(), code)
}
