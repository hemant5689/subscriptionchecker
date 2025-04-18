package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hej/internal/auth"
	"hej/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type YouTubeSubscription struct {
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ResourceId  struct {
			ChannelId string `json:"channelId"`
		} `json:"resourceId"`
	} `json:"snippet"`
}

type YouTubeSubscriptionResponse struct {
	Items []YouTubeSubscription `json:"items"`
}

// YouTubeAuthService handles YouTube authentication
type YouTubeAuthService struct {
	*auth.OAuthService
	cfg *config.Config
}

// NewYouTubeAuthService creates a new YouTube auth service
func NewYouTubeAuthService(cfg *config.Config) *YouTubeAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.YouTubeClientID,
		ClientSecret: cfg.YouTubeClientSecret,
		RedirectURL:  cfg.YouTubeRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/youtube.readonly"},
		Endpoint:     google.Endpoint,
	}
	return &YouTubeAuthService{
		OAuthService: auth.NewOAuthService(oauthConfig),
		cfg:          cfg,
	}
}

// ExchangeToken exchanges an auth code for a token
func (s *YouTubeAuthService) ExchangeToken(code string) (string, error) {
	token, err := s.ExchangeCode(code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// YouTubeService represents a YouTube API service
type YouTubeService struct {
	httpClient *http.Client
	channelID  string
}

// NewYouTubeService creates a new YouTube service with token
func NewYouTubeService(token string, cfg *config.Config) *YouTubeService {
	return &YouTubeService{
		httpClient: &http.Client{
			Transport: &tokenTransport{token: token},
		},
		channelID: cfg.YouTubeChannelID,
	}
}

// IsFollower checks if the authenticated user is subscribed to the configured channel
func (s *YouTubeService) IsFollower(_ string) (bool, error) {
	isSubscribed, _, err := s.GetSubscriptionStatus(s.channelID)
	return isSubscribed, err
}

// GetSubscriptionStatus checks if the user is subscribed to a specific channel
// and optionally returns all subscriptions
func (s *YouTubeService) GetSubscriptionStatus(channelID string) (bool, []YouTubeSubscription, error) {
	url := "https://youtube.googleapis.com/youtube/v3/subscriptions?part=snippet&mine=true&maxResults=50"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, nil, fmt.Errorf("API error: %s, %s", resp.Status, string(body))
	}

	var response YouTubeSubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for specific channel subscription if channelID was provided
	if channelID != "" {
		for _, sub := range response.Items {
			if sub.Snippet.ResourceId.ChannelId == channelID {
				return true, response.Items, nil
			}
		}
		return false, response.Items, nil
	}

	return false, response.Items, nil
}

// tokenTransport is an http.RoundTripper that adds the token to requests
type tokenTransport struct {
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
