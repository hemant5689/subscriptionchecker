package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hej/internal/auth"
	"hej/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// FacebookAuthService handles Facebook authentication
type FacebookAuthService struct {
	*auth.OAuthService
	cfg *config.Config
}

// NewFacebookAuthService creates a new Facebook auth service
func NewFacebookAuthService(cfg *config.Config) *FacebookAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.MetaAppID,
		ClientSecret: cfg.MetaAppSecret,
		RedirectURL:  cfg.MetaRedirectURI,
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	return &FacebookAuthService{
		OAuthService: auth.NewOAuthService(oauthConfig),
		cfg:          cfg,
	}
}

// ExchangeToken exchanges an auth code for a token
func (s *FacebookAuthService) ExchangeToken(code string) (string, error) {
	token, err := s.ExchangeCode(code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// FacebookUser represents a Facebook user profile
type FacebookUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FacebookFollowingResponse represents the response from the Facebook following API
type FacebookFollowingResponse struct {
	Data []FacebookUser `json:"data"`
}

// FacebookService represents a Facebook API service
type FacebookService struct {
	accessToken string
	httpClient  *http.Client
}

// NewFacebookService creates a new Facebook service with token
func NewFacebookService(token string, _ *config.Config) *FacebookService {
	return &FacebookService{
		accessToken: token,
		httpClient:  &http.Client{},
	}
}

// IsFollower checks if a user follows a page or profile
func (s *FacebookService) IsFollower(targetUserID string) (bool, error) {
	if targetUserID == "" {
		return false, fmt.Errorf("targetUserID is required")
	}

	// Get the current user's profile
	me, err := s.getProfile()
	if err != nil {
		return false, err
	}

	// Check if the user follows the target
	isFollowing, err := s.checkFollowing(me.ID, targetUserID)
	if err != nil {
		return false, err
	}

	return isFollowing, nil
}

// getProfile gets the current user's profile
func (s *FacebookService) getProfile() (*FacebookUser, error) {
	reqURL := fmt.Sprintf("https://graph.facebook.com/v18.0/me?access_token=%s", s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var user FacebookUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// checkFollowing checks if a user follows another user
func (s *FacebookService) checkFollowing(userID, targetPageID string) (bool, error) {
	reqURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/subscribedto?access_token=%s", userID, s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return false, fmt.Errorf("failed to check following: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var response FacebookFollowingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the target user is in the list of friends
	for _, user := range response.Data {
		if user.ID == targetPageID {
			return true, nil
		}
	}

	return false, nil
}
