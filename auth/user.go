package auth

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
)

type contextKey string

// use module path to make it unique
const userCtxKey contextKey = "github.com/irvankadhafi/user-balance-transfer-service/auth.User"

// SetUserToCtx set user to context
func SetUserToCtx(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// GetUserFromCtx get user from context
func GetUserFromCtx(ctx context.Context) *User {
	user, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return nil
	}
	return &user
}

// User represent an authenticated user
type User struct {
	ID        int `json:"id"`
	SessionID int `json:"session_id"`
}

// NewUserFromSession return new user from session
func NewUserFromSession(sess model.Session) User {
	return User{
		ID:        sess.UserID,
		SessionID: sess.ID,
	}
}
