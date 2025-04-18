package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hej/internal/auth"
	"hej/internal/config"

	"golang.org/x/oauth2"
)

// Discord OAuth endpoints
var discordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

// DiscordAuthService handles Discord authentication
type DiscordAuthService struct {
	*auth.OAuthService
	cfg *config.Config
}

// NewDiscordAuthService creates a new Discord auth service
func NewDiscordAuthService(cfg *config.Config) *DiscordAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.DiscordClientID,
		ClientSecret: cfg.DiscordClientSecret,
		RedirectURL:  cfg.DiscordRedirectURI,
		Scopes:       []string{"identify", "guilds"},
		Endpoint:     discordEndpoint,
	}
	return &DiscordAuthService{
		OAuthService: auth.NewOAuthService(oauthConfig),
		cfg:          cfg,
	}
}

// ExchangeToken exchanges an auth code for a token
func (s *DiscordAuthService) ExchangeToken(code string) (string, error) {
	token, err := s.ExchangeCode(code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// DiscordUser represents a Discord user profile
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

// DiscordGuild represents a Discord server/guild
type DiscordGuild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DiscordService represents a Discord API service
type DiscordService struct {
	accessToken string
	httpClient  *http.Client
	serverID    string
}

// NewDiscordService creates a new Discord service with token
func NewDiscordService(token string, cfg *config.Config) *DiscordService {
	return &DiscordService{
		accessToken: token,
		httpClient:  &http.Client{},
		serverID:    cfg.DiscordServerID,
	}
}

// IsFollower checks if a user is a member of the specified Discord server
func (s *DiscordService) IsFollower(_ string) (bool, error) {
	if s.serverID == "" {
		return false, fmt.Errorf("server ID is not configured")
	}

	// Get the user's guilds (servers)
	guilds, err := s.getUserGuilds()
	if err != nil {
		return false, err
	}

	// Check if the user is a member of the target server
	for _, guild := range guilds {
		if guild.ID == s.serverID {
			return true, nil
		}
	}

	return false, nil
}

// getUserProfile gets the authenticated user's Discord profile
func (s *DiscordService) getUserProfile() (*DiscordUser, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
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
		return nil, fmt.Errorf("Discord API error: %s, %s", resp.Status, string(body))
	}

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// getUserGuilds gets the servers/guilds the user is a member of
func (s *DiscordService) getUserGuilds() ([]DiscordGuild, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me/guilds", nil)
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
		return nil, fmt.Errorf("Discord API error: %s, %s", resp.Status, string(body))
	}

	var guilds []DiscordGuild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return guilds, nil
}
