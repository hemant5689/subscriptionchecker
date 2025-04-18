package handler

import (
	"hej/internal/config"
	"hej/internal/platform"
	"hej/pkg/utils"
	"net/http"
)

// FacebookHandler handles Facebook-related requests
type FacebookHandler struct {
	cfg *config.Config
}

// NewFacebookHandler creates a new Facebook handler
func NewFacebookHandler(cfg *config.Config) *FacebookHandler {
	return &FacebookHandler{
		cfg: cfg,
	}
}

// Login handles the Facebook auth login request
func (h *FacebookHandler) Login(w http.ResponseWriter, r *http.Request) {
	authService := platform.NewFacebookAuthService(h.cfg)
	http.Redirect(w, r, authService.GetAuthURL(), http.StatusTemporaryRedirect)
}

// Callback handles the Facebook auth callback
func (h *FacebookHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Code not found")
		return
	}

	authService := platform.NewFacebookAuthService(h.cfg)
	token, err := authService.ExchangeToken(code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

// CheckFollower checks if a user follows a Facebook profile
func (h *FacebookHandler) CheckFollower(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	targetID := r.URL.Query().Get("targetId")
	if targetID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Target ID is required")
		return
	}

	service := platform.NewFacebookService(token, h.cfg)
	isFollowing, err := service.IsFollower(targetID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check follower status: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"isFollowing": isFollowing,
	})
}
