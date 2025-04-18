package handler

import (
	"hej/internal/config"
	"hej/internal/platform"
	"hej/pkg/utils"
	"net/http"
)

// YouTubeHandler handles YouTube-related requests
type YouTubeHandler struct {
	cfg *config.Config
}

// NewYouTubeHandler creates a new YouTube handler
func NewYouTubeHandler(cfg *config.Config) *YouTubeHandler {
	return &YouTubeHandler{
		cfg: cfg,
	}
}

// Login handles the YouTube auth login request
func (h *YouTubeHandler) Login(w http.ResponseWriter, r *http.Request) {
	authService := platform.NewYouTubeAuthService(h.cfg)
	http.Redirect(w, r, authService.GetAuthURL(), http.StatusTemporaryRedirect)
}

// Callback handles the YouTube auth callback
func (h *YouTubeHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Code not found")
		return
	}

	authService := platform.NewYouTubeAuthService(h.cfg)
	token, err := authService.ExchangeToken(code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

// CheckSubscription checks if a user is subscribed to a YouTube channel
func (h *YouTubeHandler) CheckSubscription(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	service := platform.NewYouTubeService(token, h.cfg)
	isSubscribed, err := service.IsFollower("")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check subscription: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"isSubscribed": isSubscribed,
	})
}
