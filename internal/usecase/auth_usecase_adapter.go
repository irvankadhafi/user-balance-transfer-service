package usecase

import (
	"context"
	"errors"
	"github.com/irvankadhafi/user-balance-transfer-service/auth"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserAutherAdapter adapter for auth.UserAuthenticator
type UserAutherAdapter struct {
	authUsecase model.AuthUsecase
}

// NewUserAutherAdapter constructor
func NewUserAutherAdapter(authUsecase model.AuthUsecase) *UserAutherAdapter {
	return &UserAutherAdapter{
		authUsecase: authUsecase,
	}
}

// AuthenticateToken authenticate access token
func (a *UserAutherAdapter) AuthenticateToken(ctx context.Context, accessToken string) (*auth.User, error) {
	user, err := a.authUsecase.AuthenticateToken(ctx, accessToken)
	if errors.Is(err, ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, ErrAccessTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return newAuthUser(user), nil
}

func newAuthUser(user *model.User) *auth.User {
	if user == nil {
		return nil
	}
	return &auth.User{
		ID:        user.ID,
		SessionID: user.SessionID,
	}
}
