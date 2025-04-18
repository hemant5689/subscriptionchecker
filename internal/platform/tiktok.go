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

// TiktokAuthService handles Tiktok authentication
type TiktokAuthService struct {
	*auth.OAuthService
	cfg *config.Config
}

// NewTiktokAuthService creates a new Tiktok auth service
func NewTiktokAuthService(cfg *config.Config) *TiktokAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.MetaAppID,
		ClientSecret: cfg.MetaAppSecret,
		RedirectURL:  cfg.MetaRedirectURI,
		Scopes:       []string{"user_profile", "instagram_basic"},
		Endpoint:     facebook.Endpoint,
	}
	return &TiktokAuthService{
		OAuthService: auth.NewOAuthService(oauthConfig),
		cfg:          cfg,
	}
}

// ExchangeToken exchanges an auth code for a token
func (s *TiktokAuthService) ExchangeToken(code string) (string, error) {
	token, err := s.ExchangeCode(code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// TiktokProfile represents a Tiktok user profile
type TiktokProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// TiktokFollowerResponse represents the response from the Tiktok followers API
type TiktokFollowerResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

// TiktokService represents a Tiktok API service
type TiktokService struct {
	accessToken string
	httpClient  *http.Client
}

// NewInstagramService creates a new Instagram service with token
func NewTiktokService(token string, _ *config.Config) *TiktokService {
	return &TiktokService{
		accessToken: token,
		httpClient:  &http.Client{},
	}
}

// IsFollower checks if a user follows another Tiktok user
func (s *TiktokService) IsFollower(targetUsername string) (bool, error) {
	if targetUsername == "" {
		return false, fmt.Errorf("targetUsername is required")
	}

	// First, get the user's Tiktok profile ID
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
func (s *TiktokService) getProfile() (*TiktokProfile, error) {
	// Get the Tiktok user ID from the Tiktok API
	reqURL := fmt.Sprintf("https://api.tiktok.com/v2/user/info?access_token=%s", s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var profile TiktokProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &profile, nil
}

// getUserByUsername gets a user profile by username
func (s *TiktokService) getUserByUsername(username string) (*TiktokProfile, error) {
	// Note: This is a simplified implementation and might need adjustment based on Tiktok's API
	reqURL := fmt.Sprintf("https://api.tiktok.com/v2/user/info?access_token=%s", s.accessToken)

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

	return &TiktokProfile{
		ID:       data.AuthorID,
		Username: data.AuthorName,
	}, nil
}

// checkFollowing checks if a user follows another user
func (s *TiktokService) checkFollowing(userID, targetUserID string) (bool, error) {
	// Tiktok API to check followers
	reqURL := fmt.Sprintf("https://api.tiktok.com/v2/user/following?access_token=%s", s.accessToken)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return false, fmt.Errorf("failed to check following: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var response TiktokFollowerResponse
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
