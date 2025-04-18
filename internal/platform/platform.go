package platform

// FollowerChecker defines the interface for checking if a user follows another user
type FollowerChecker interface {
	// IsFollower checks if a user with the given username follows the authenticated user
	IsFollower(username string) (bool, error)
}

// AuthService defines common authentication methods
type AuthService interface {
	// GetAuthURL returns the OAuth URL for authentication
	GetAuthURL() string

	// ExchangeToken exchanges an authorization code for an access token
	ExchangeToken(code string) (string, error)
}
