package router

import (
	"hej/internal/config"
	"hej/internal/handler"
	"net/http"
)

// Router handles routing HTTP requests
type Router struct {
	cfg *config.Config
}

// NewRouter creates a new router instance
func NewRouter(cfg *config.Config) *Router {
	return &Router{
		cfg: cfg,
	}
}

// Setup configures all application routes
func (r *Router) Setup() {
	// Initialize handlers
	youtubeHandler := handler.NewYouTubeHandler(r.cfg)
	facebookHandler := handler.NewFacebookHandler(r.cfg)
	instagramHandler := handler.NewInstagramHandler(r.cfg)
	discordHandler := handler.NewDiscordHandler(r.cfg)
	twitterHandler := handler.NewTwitterHandler(r.cfg)

	// YouTube routes
	http.HandleFunc("/auth/youtube/login", youtubeHandler.Login)
	http.HandleFunc("/auth/youtube/callback", youtubeHandler.Callback)
	http.HandleFunc("/check-youtube-subscription", youtubeHandler.CheckSubscription)

	// Facebook routes
	http.HandleFunc("/auth/facebook/login", facebookHandler.Login)
	http.HandleFunc("/auth/facebook/callback", facebookHandler.Callback)
	http.HandleFunc("/check-facebook-follower", facebookHandler.CheckFollower)

	// Instagram routes
	http.HandleFunc("/auth/instagram/login", instagramHandler.Login)
	http.HandleFunc("/auth/instagram/callback", instagramHandler.Callback)
	http.HandleFunc("/check-instagram-follower", instagramHandler.CheckFollower)

	// Discord routes
	http.HandleFunc("/auth/discord/login", discordHandler.Login)
	http.HandleFunc("/auth/discord/callback", discordHandler.Callback)
	http.HandleFunc("/check-discord-server", discordHandler.CheckServerMembership)

	// Twitter routes
	http.HandleFunc("/auth/twitter/login", twitterHandler.Login)
	http.HandleFunc("/auth/twitter/callback", twitterHandler.Callback)
	http.HandleFunc("/check-twitter-follower", twitterHandler.CheckFollower)

	// Add more platform routes as needed:
	// - Discord
	// - Twitter
}
