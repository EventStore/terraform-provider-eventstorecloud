package client

type tokenData struct {
	AccessToken  accessToken `json:"access_token"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	Scope        string      `json:"scope"`
	ExpiresIn    int         `json:"expires_in"`
	TokenType    string      `json:"token_type"`
}
