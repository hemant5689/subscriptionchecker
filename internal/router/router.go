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
	tiktokHandler := handler.NewTiktokHandler(r.cfg)
	// YouTube routes
	http.HandleFunc("/youtube/login", youtubeHandler.Login)
	http.HandleFunc("/youtube/callback", youtubeHandler.Callback)
	http.HandleFunc("/youtube/check-subscription", youtubeHandler.CheckSubscription)

	// Facebook routes
	http.HandleFunc("/facebook/login", facebookHandler.Login)
	http.HandleFunc("/facebook/callback", facebookHandler.Callback)
	http.HandleFunc("/facebook/check-follower", facebookHandler.CheckFollower)

	// Instagram routes
	http.HandleFunc("/instagram/login", instagramHandler.Login)
	http.HandleFunc("/instagram/callback", instagramHandler.Callback)
	http.HandleFunc("/instagram/check-follower", instagramHandler.CheckFollower)

	// Discord routes
	http.HandleFunc("/discord/login", discordHandler.Login)
	http.HandleFunc("/discord/callback", discordHandler.Callback)
	http.HandleFunc("/discord/check-server", discordHandler.CheckServerMembership)

	// Twitter routes
	http.HandleFunc("/twitter/login", twitterHandler.Login)
	http.HandleFunc("/twitter/callback", twitterHandler.Callback)
	http.HandleFunc("/twitter/check-follower", twitterHandler.CheckFollower)

	//tiktok routes
	http.HandleFunc("/tiktok/login", tiktokHandler.Login)
	http.HandleFunc("/tiktok/callback", tiktokHandler.Callback)
	http.HandleFunc("/tiktok/check-follower", tiktokHandler.CheckFollower)
}
