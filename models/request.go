package models

type GenerateTokenRequest struct {
	Username string `json:"username"`
	DeviceID string `json:"device_id"`
	Email    string `json:"email"`
}

// Response structure for token generation
type TokenResponse struct {
	UserName     string `json:"user_id"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
