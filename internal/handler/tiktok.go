package handler

import (
	"hej/internal/config"
	"hej/internal/platform"
	"hej/pkg/utils"
	"net/http"
)

// TiktokHandler handles Tiktok-related requests
type TiktokHandler struct {
	cfg *config.Config
}

// NewTiktokHandler creates a new Tiktok handler
func NewTiktokHandler(cfg *config.Config) *TiktokHandler {
	return &TiktokHandler{
		cfg: cfg,
	}
}

// Login handles the Tiktok auth login request
func (h *TiktokHandler) Login(w http.ResponseWriter, r *http.Request) {
	authService := platform.NewTiktokAuthService(h.cfg)
	http.Redirect(w, r, authService.GetAuthURL(), http.StatusTemporaryRedirect)
}

// Callback handles the Tiktok auth callback
func (h *TiktokHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Code not found")
		return
	}

	authService := platform.NewTiktokAuthService(h.cfg)
	token, err := authService.ExchangeToken(code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

// CheckFollower checks if a user is a follower of a specific Tiktok user
func (h *TiktokHandler) CheckFollower(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	service := platform.NewTiktokService(token, h.cfg)
	isFollower, err := service.IsFollower("")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check follower: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"isFollower": isFollower,
	})
}
