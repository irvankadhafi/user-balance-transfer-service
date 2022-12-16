package model

import "context"

// LoginRequest request
type LoginRequest struct {
	Username      string `json:"username"`
	PlainPassword string `json:"plain_password"`
	IPAddress     string `json:"ip_address"`
	UserAgent     string `json:"user_agent"`
}

// RefreshTokenRequest request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
}

// AuthUsecase usecases about IAM
type AuthUsecase interface {
	LoginByUsernamePassword(ctx context.Context, req LoginRequest) (*Session, error)

	// AuthenticateToken authenticate the given token
	AuthenticateToken(ctx context.Context, accessToken string) (*User, error)

	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*Session, error)
	DeleteSessionByID(ctx context.Context, sessionID int) error
}
