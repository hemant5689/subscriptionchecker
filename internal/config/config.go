package config

import (
	"os"
	"sync"
)

// Config holds all application configuration
type Config struct {
	Port string

	// YouTube
	YouTubeClientID     string
	YouTubeClientSecret string
	YouTubeRedirectURL  string
	YouTubeChannelID    string

	// Facebook/Instagram (Meta)
	MetaAppID       string
	MetaAppSecret   string
	MetaRedirectURI string

	// Discord
	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURI  string
	DiscordServerID     string

	// Twitter
	TwitterClientID     string
	TwitterClientSecret string
	TwitterRedirectURI  string

	// Tiktok
	TiktokClientID     string
	TiktokClientSecret string
	TiktokRedirectURI  string
}

var (
	config Config
	once   sync.Once
)

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	once.Do(func() {
		config = Config{
			Port: getEnvOrDefault("PORT", "8080"),

			// YouTube
			YouTubeClientID:     os.Getenv("YT_CLIENT_ID"),
			YouTubeClientSecret: os.Getenv("YT_CLIENT_SECRET"),
			YouTubeRedirectURL:  os.Getenv("YT_REDIRECT_URL"),
			YouTubeChannelID:    os.Getenv("YT_CHANNEL_ID"),

			// Facebook/Instagram (Meta)
			MetaAppID:       os.Getenv("META_APP_ID"),
			MetaAppSecret:   os.Getenv("META_APP_SECRET"),
			MetaRedirectURI: os.Getenv("META_REDIRECT_URI"),

			// Discord
			DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			DiscordRedirectURI:  os.Getenv("DISCORD_REDIRECT_URI"),
			DiscordServerID:     os.Getenv("DISCORD_SERVER_ID"),

			// Twitter
			TwitterClientID:     os.Getenv("TWITTER_CLIENT_ID"),
			TwitterClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
			TwitterRedirectURI:  os.Getenv("TWITTER_REDIRECT_URI"),

			// Tiktok
			TiktokClientID:     os.Getenv("TIKTOK_CLIENT_ID"),
			TiktokClientSecret: os.Getenv("TIKTOK_CLIENT_SECRET"),
			TiktokRedirectURI:  os.Getenv("TIKTOK_REDIRECT_URI"),
		}
	})

	return &config
}

// getEnvOrDefault returns the environment variable value or a default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
