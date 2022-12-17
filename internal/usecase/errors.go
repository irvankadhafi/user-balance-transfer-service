package usecase

import "errors"

// errors ..
var (
	ErrNotFound            = errors.New("not found")
	ErrAccessTokenExpired  = errors.New("access token expired")
	ErrFailedPrecondition  = errors.New("precondition failed")
	ErrDuplicateUser       = errors.New("user already exist")
	ErrLoginMaxAttempts    = errors.New("user is locked from logging")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrBalanceNotEnough    = errors.New("balance not enough")
)
