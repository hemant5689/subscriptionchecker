package handler

import (
	"hej/internal/config"
	"hej/internal/platform"
	"hej/pkg/utils"
	"net/http"
)

// DiscordHandler handles Discord-related requests
type DiscordHandler struct {
	cfg *config.Config
}

// NewDiscordHandler creates a new Discord handler
func NewDiscordHandler(cfg *config.Config) *DiscordHandler {
	return &DiscordHandler{
		cfg: cfg,
	}
}

// Login handles the Discord auth login request
func (h *DiscordHandler) Login(w http.ResponseWriter, r *http.Request) {
	authService := platform.NewDiscordAuthService(h.cfg)
	http.Redirect(w, r, authService.GetAuthURL(), http.StatusTemporaryRedirect)
}

// Callback handles the Discord auth callback
func (h *DiscordHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Code not found")
		return
	}

	authService := platform.NewDiscordAuthService(h.cfg)
	token, err := authService.ExchangeToken(code)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"accessToken": token,
	})
}

// CheckServerMembership checks if a user is a member of a specific Discord server
func (h *DiscordHandler) CheckServerMembership(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	service := platform.NewDiscordService(token, h.cfg)
	isMember, err := service.IsFollower("")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check server membership: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"isMember": isMember,
	})
}
