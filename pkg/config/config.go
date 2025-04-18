package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func LoadOAuthConfig() *OAuthConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &OAuthConfig{
		ClientID:     os.Getenv("YT_CLIENT_ID"),
		ClientSecret: os.Getenv("YT_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("YT_REDIRECT_URL"),
	}
}
