package model

import (
	"context"
)

// User :nodoc:
type User struct {
	ID       int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" gorm:"->:false;<-"` // gorm create & update only (disabled read from db)

	SessionID int `json:"session_id" gorm:"-"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	IsLoginByEmailPasswordLocked(ctx context.Context, email string) (bool, error)
	IncrementLoginByEmailPasswordRetryAttempts(ctx context.Context, username string) error
	FindPasswordByID(ctx context.Context, id int) ([]byte, error)
}

type UserUsecase interface {
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	FindByID(ctx context.Context, userID int) (*User, error)
}

// CreateUserInput :nodoc:
type CreateUserInput struct {
	Email                string `json:"email" validate:"required,email"`
	Username             string `json:"username" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=Password"`
}

// Validate validate user input body
func (c *CreateUserInput) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	return nil
}
