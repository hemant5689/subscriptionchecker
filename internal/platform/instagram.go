package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hej/internal/auth"
	"hej/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook" // Instagram uses the same OAuth endpoint as Facebook
)

// InstagramAuthService handles Instagram authentication
type InstagramAuthService struct {
	*auth.OAuthService
	cfg *config.Config
}

// NewInstagramAuthService creates a new Instagram auth service
func NewInstagramAuthService(cfg *config.Config) *InstagramAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.MetaAppID,
		ClientSecret: cfg.MetaAppSecret,
		RedirectURL:  cfg.MetaRedirectURI,
		Scopes:       []string{"user_profile", "instagram_basic"},
		Endpoint:     facebook.Endpoint,
	}
	return &InstagramAuthService{
		OAuthService: auth.NewOAuthService(oauthConfig),
		cfg:          cfg,
	}
}

// ExchangeToken exchanges an auth code for a token
func (s *InstagramAuthService) ExchangeToken(code string) (string, error) {
	token, err := s.ExchangeCode(code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// InstagramProfile represents an Instagram user profile
type InstagramProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// InstagramFollowerResponse represents the response from the Instagram followers API
type InstagramFollowerResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

// InstagramService represents an Instagram API service
type InstagramService struct {
	accessToken string
	httpClient  *http.Client
}

// NewInstagramService creates a new Instagram service with token
func NewInstagramService(token string, _ *config.Config) *InstagramService {
	return &InstagramService{
		accessToken: token,
		httpClient:  &http.Client{},
	}
}

// IsFollower checks if a user follows another Instagram user
func (s *InstagramService) IsFollower(targetUsername string) (bool, error) {
	if targetUsername == "" {
		return false, fmt.Errorf("targetUsername is required")
	}

	// First, get the user's Instagram profile ID
	profile, err := s.getProfile()
	if err != nil {
		return false, err
	}

	// Then check if the target user is in the followers list
	targetProfile, err := s.getUserByUsername(targetUsername)
	if err != nil {
		return false, err
	}

	isFollowing, err := s.checkFollowing(profile.ID, targetProfile.ID)
	if err != nil {
		return false, err
	}

	return isFollowing, nil
}

// getProfile gets the authenticated user's Instagram profile
func (s *InstagramService) getProfile() (*InstagramProfile, error) {
	// Get the Instagram user ID from the Facebook Graph API
	reqURL := fmt.Sprintf("https://graph.facebook.com/v18.0/me?fields=id,name&access_token=%s", s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var profile InstagramProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &profile, nil
}

// getUserByUsername gets a user profile by username
func (s *InstagramService) getUserByUsername(username string) (*InstagramProfile, error) {
	// Note: This is a simplified implementation and might need adjustment based on Instagram's API
	reqURL := fmt.Sprintf("https://graph.facebook.com/v18.0/instagram_oembed?url=https://www.instagram.com/%s/&access_token=%s", username, s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var data struct {
		AuthorName string `json:"author_name"`
		AuthorID   string `json:"author_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &InstagramProfile{
		ID:       data.AuthorID,
		Username: data.AuthorName,
	}, nil
}

// checkFollowing checks if a user follows another user
func (s *InstagramService) checkFollowing(userID, targetUserID string) (bool, error) {
	// Instagram API to check followers
	reqURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/following?access_token=%s", userID, s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return false, fmt.Errorf("failed to check following: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var response InstagramFollowerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if target user ID is in the list
	for _, user := range response.Data {
		if user.ID == targetUserID {
			return true, nil
		}
	}

	return false, nil
}
