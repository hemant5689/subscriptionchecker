package handler

import (
	"hej/internal/config"
	"hej/internal/platform"
	"hej/pkg/utils"
	"net/http"
)

// TwitterHandler handles Twitter-related requests
type TwitterHandler struct {
	cfg *config.Config
}

// NewTwitterHandler creates a new Twitter handler
func NewTwitterHandler(cfg *config.Config) *TwitterHandler {
	return &TwitterHandler{
		cfg: cfg,
	}
}

// Login handles the Twitter auth login request
func (h *TwitterHandler) Login(w http.ResponseWriter, r *http.Request) {
	authService := platform.NewTwitterAuthService(h.cfg)
	http.Redirect(w, r, authService.GetAuthURL(), http.StatusTemporaryRedirect)
}

// Callback handles the Twitter auth callback
func (h *TwitterHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Code not found")
		return
	}

	authService := platform.NewTwitterAuthService(h.cfg)
	token, err := authService.ExchangeToken(code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

// CheckFollower checks if a user follows a Twitter account
func (h *TwitterHandler) CheckFollower(w http.ResponseWriter, r *http.Request) {
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

	service := platform.NewTwitterService(token, h.cfg)
	isFollowing, err := service.IsFollower(targetUsername)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check follower status: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"isFollowing": isFollowing,
	})
}
