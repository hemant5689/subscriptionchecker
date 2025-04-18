package handler

import (
	"hej/internal/config"
	"hej/internal/platform"
	"hej/pkg/utils"
	"net/http"
)

// InstagramHandler handles Instagram-related requests
type InstagramHandler struct {
	cfg *config.Config
}

// NewInstagramHandler creates a new Instagram handler
func NewInstagramHandler(cfg *config.Config) *InstagramHandler {
	return &InstagramHandler{
		cfg: cfg,
	}
}

// Login handles the Instagram auth login request
func (h *InstagramHandler) Login(w http.ResponseWriter, r *http.Request) {
	authService := platform.NewInstagramAuthService(h.cfg)
	http.Redirect(w, r, authService.GetAuthURL(), http.StatusTemporaryRedirect)
}

// Callback handles the Instagram auth callback
func (h *InstagramHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Code not found")
		return
	}

	authService := platform.NewInstagramAuthService(h.cfg)
	token, err := authService.ExchangeToken(code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

// CheckFollower checks if a user follows an Instagram account
func (h *InstagramHandler) CheckFollower(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	targetUsername := r.URL.Query().Get("username")
	if targetUsername == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Target username is required")
		return
	}

	service := platform.NewInstagramService(token, h.cfg)
	isFollowing, err := service.IsFollower(targetUsername)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check follower status: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"isFollowing": isFollowing,
	})
}
