package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"hej/internal/auth"
	"hej/internal/config"

	"golang.org/x/oauth2"
)

// Twitter OAuth endpoints
var twitterEndpoint = oauth2.Endpoint{
	AuthURL:  "https://twitter.com/i/oauth2/authorize",
	TokenURL: "https://api.twitter.com/2/oauth2/token",
}

// TwitterAuthService handles Twitter authentication
type TwitterAuthService struct {
	*auth.OAuthService
	cfg *config.Config
}

// NewTwitterAuthService creates a new Twitter auth service
func NewTwitterAuthService(cfg *config.Config) *TwitterAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.TwitterClientID,
		ClientSecret: cfg.TwitterClientSecret,
		RedirectURL:  cfg.TwitterRedirectURI,
		Scopes:       []string{"tweet.read", "users.read", "follows.read"},
		Endpoint:     twitterEndpoint,
	}
	return &TwitterAuthService{
		OAuthService: auth.NewOAuthService(oauthConfig),
		cfg:          cfg,
	}
}

// ExchangeToken exchanges an auth code for a token
func (s *TwitterAuthService) ExchangeToken(code string) (string, error) {
	token, err := s.ExchangeCode(code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// TwitterUser represents a Twitter user profile
type TwitterUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// TwitterUsersResponse represents the response from the Twitter users API
type TwitterUsersResponse struct {
	Data TwitterUser `json:"data"`
}

// TwitterFollowingResponse represents the response from the Twitter following API
type TwitterFollowingResponse struct {
	Data []TwitterUser `json:"data"`
	Meta struct {
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
	} `json:"meta"`
}

// TwitterService represents a Twitter API service
type TwitterService struct {
	accessToken string
	httpClient  *http.Client
}

// NewTwitterService creates a new Twitter service with token
func NewTwitterService(token string, _ *config.Config) *TwitterService {
	return &TwitterService{
		accessToken: token,
		httpClient:  &http.Client{},
	}
}

// IsFollower checks if a user follows a Twitter account
func (s *TwitterService) IsFollower(targetUsername string) (bool, error) {
	if targetUsername == "" {
		return false, fmt.Errorf("target username is required")
	}

	// First, get the user's profile
	me, err := s.getProfile()
	if err != nil {
		return false, err
	}

	// Then get the target user's profile
	targetUser, err := s.getUserByUsername(targetUsername)
	if err != nil {
		return false, err
	}

	// Check if the user follows the target
	isFollowing, err := s.checkFollowing(me.ID, targetUser.ID)
	if err != nil {
		return false, err
	}

	return isFollowing, nil
}

// getProfile gets the authenticated user's Twitter profile
func (s *TwitterService) getProfile() (*TwitterUser, error) {
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Twitter API error: %s, %s", resp.Status, string(body))
	}

	var response TwitterUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// getUserByUsername gets a Twitter user by username
func (s *TwitterService) getUserByUsername(username string) (*TwitterUser, error) {
	endpoint := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", url.PathEscape(username))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Twitter API error: %s, %s", resp.Status, string(body))
	}

	var response TwitterUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// checkFollowing checks if a user follows another user
func (s *TwitterService) checkFollowing(userID, targetUserID string) (bool, error) {
	endpoint := fmt.Sprintf("https://api.twitter.com/2/users/%s/following", userID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	// Set query parameters
	q := req.URL.Query()
	q.Add("max_results", "100")
	req.URL.RawQuery = q.Encode()

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("Twitter API error: %s, %s", resp.Status, string(body))
	}

	var response TwitterFollowingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the target user is in the list of followed accounts
	for _, user := range response.Data {
		if user.ID == targetUserID {
			return true, nil
		}
	}

	// TODO: Handle pagination if there are more than max_results followers

	return false, nil
}
