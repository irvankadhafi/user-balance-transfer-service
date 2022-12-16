package usecase

import (
	"context"
	"errors"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/config"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/helper"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type authUsecase struct {
	userUsecase model.UserUsecase
	userRepo    model.UserRepository
	sessionRepo model.SessionRepository
}

// NewAuthUsecase :nodoc:
func NewAuthUsecase(
	userRepo model.UserRepository,
	sessionRepo model.SessionRepository,
	userUsecase model.UserUsecase,
) model.AuthUsecase {
	return &authUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		userUsecase: userUsecase,
	}
}

// LoginByUsernamePassword login the user by username & password
func (a *authUsecase) LoginByUsernamePassword(ctx context.Context, req model.LoginRequest) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"username":  req.Username,
		"ip":        req.IPAddress,
		"userAgent": req.UserAgent,
	})

	isLocked, err := a.userRepo.IsLoginByUsernamePasswordLocked(ctx, req.Username)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if isLocked {
		return nil, ErrLoginMaxAttempts
	}

	user, err := a.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	cipherPass, err := a.userRepo.FindPasswordByID(ctx, user.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if cipherPass == nil {
		logger.Error(err)
		return nil, errors.New("unexpected: no password found")
	}

	if !helper.IsHashedStringMatch([]byte(req.PlainPassword), cipherPass) {
		// obscure the error if the password does not match
		if err := a.userRepo.IncrementLoginByUsernamePasswordRetryAttempts(ctx, req.Username); err != nil {
			logger.Error(err)
			return nil, err
		}

		return nil, ErrUnauthorized
	}

	logger = logger.WithField("userID", user.ID)
	accessToken, err := generateToken(a.sessionRepo, user.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	refreshToken, err := generateToken(a.sessionRepo, user.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &model.Session{
		UserID:                user.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  now.Add(config.AccessTokenDuration()),
		RefreshTokenExpiredAt: now.Add(config.RefreshTokenDuration()),
		IPAddress:             req.IPAddress,
		UserAgent:             req.UserAgent,
		Location:              "-",
	}

	if err = a.sessionRepo.Create(ctx, session); err != nil {
		logger.Error(err)
		return nil, err
	}
	return session, nil
}

// AuthenticateToken authenticate the given access token and return the corresponding user
func (a *authUsecase) AuthenticateToken(ctx context.Context, accessToken string) (*model.User, error) {
	session, err := a.sessionRepo.FindByToken(ctx, model.AccessToken, accessToken)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if session == nil {
		return nil, ErrNotFound
	}

	if session.IsAccessTokenExpired() {
		return nil, ErrAccessTokenExpired
	}

	user, err := a.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		logrus.WithField("userID", session.UserID).Error(err)
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	user.SessionID = session.ID

	return user, nil
}

// RefreshToken refresh the user's access and refresh token
func (a *authUsecase) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":                 utils.DumpIncomingContext(ctx),
		"refreshTokenRequest": utils.Dump(req),
	})

	session, err := a.sessionRepo.FindByToken(ctx, model.RefreshToken, req.RefreshToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if session == nil {
		logger.Error(ErrNotFound)
		return nil, ErrNotFound
	}

	user, err := a.userRepo.FindByID(ctx, session.UserID)
	switch {
	case err != nil:
		logger.WithField("userID", session.UserID).Error(err)
		return nil, err
	case user == nil:
		logger.WithField("userID", session.UserID).Error(ErrNotFound)
		return nil, ErrNotFound
	}

	// old session is used to delete the old session cache
	oldSess := *session

	if session.RefreshTokenExpiredAt.Before(time.Now()) {
		logger.Error(ErrRefreshTokenExpired)
		return nil, ErrRefreshTokenExpired
	}

	newAccessToken, err := generateToken(a.sessionRepo, session.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newRefreshToken, err := generateToken(a.sessionRepo, session.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	session.AccessToken = newAccessToken
	session.RefreshToken = newRefreshToken
	session.IPAddress = req.IPAddress
	session.UserAgent = req.UserAgent

	now := time.Now()
	session.AccessTokenExpiredAt = now.Add(config.AccessTokenDuration())
	session.RefreshTokenExpiredAt = now.Add(config.RefreshTokenDuration())

	session, err = a.sessionRepo.RefreshToken(ctx, &oldSess, session)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return session, nil
}

// DeleteSessionByID deletes session by id.
func (a *authUsecase) DeleteSessionByID(ctx context.Context, sessionID int) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"sessionID": utils.Dump(sessionID),
	})

	session, err := a.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if session == nil {
		return ErrNotFound
	}

	err = a.sessionRepo.Delete(ctx, session)
	if err != nil {
		logger.Error(err)
	}

	return err
}
