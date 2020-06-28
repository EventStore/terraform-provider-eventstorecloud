package client

type tokenData struct {
	AccessToken  accessToken `json:"access_token"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	Scope        string      `json:"scope"`
	ExpiresIn    int         `json:"expires_in"`
	TokenType    string      `json:"token_type"`
}

// Check if token data has a refresh token
func (t *tokenData) hasRefresh() bool {
	if t.RefreshToken != "" {
		return true
	}

	return false
}
