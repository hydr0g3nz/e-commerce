package domain

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // Only included in login response
	ExpiresIn    int64  `json:"expires_in"`              // Access token expiration in seconds
}
type TokenMetadata struct {
	UserId string `json:"user_id"` // User ID associated with the token
	Uuid   string `json:"uuid"`    // Unique token identifier
	Exp    int64  `json:"exp"`     // Token expiration timestamp in seconds
}
type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
