package httpsvc

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/auth"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/usecase"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type loginResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

func (s *Service) handleLoginByUsernamePassword() echo.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(c echo.Context) error {
		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		session, err := s.authUsecase.LoginByUsernamePassword(c.Request().Context(), model.LoginRequest{
			Username:      req.Username,
			PlainPassword: req.Password,
			IPAddress:     c.RealIP(),
			UserAgent:     c.Request().UserAgent(),
		})
		switch err {
		case nil:
			break
		case usecase.ErrNotFound, usecase.ErrUnauthorized:
			return ErrEmailPasswordNotMatch
		case usecase.ErrLoginMaxAttempts:
			return ErrLoginByEmailPasswordLocked
		default:
			logrus.Error(err)
			return ErrInternal
		}

		res := loginResponse{
			TokenType:             "Bearer",
			AccessToken:           session.AccessToken,
			AccessTokenExpiresAt:  utils.FormatTimeRFC3339(&session.AccessTokenExpiredAt),
			RefreshToken:          session.RefreshToken,
			RefreshTokenExpiresAt: utils.FormatTimeRFC3339(&session.RefreshTokenExpiredAt),
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (s *Service) handleRefreshToken() echo.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(c echo.Context) error {
		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		session, err := s.authUsecase.RefreshToken(c.Request().Context(), model.RefreshTokenRequest{
			RefreshToken: req.RefreshToken,
			IPAddress:    c.RealIP(),
			UserAgent:    c.Request().UserAgent(),
		})
		switch err {
		case nil:
		//case usecase.ErrRefreshTokenExpired, usecase.ErrNotFound:
		//	return ErrUnauthenticated
		//case usecase.ErrDiscrepantAppID:
		//	return ErrUnauthorized
		default:
			logrus.Error(err)
			return ErrInternal
		}

		res := loginResponse{
			AccessToken:           session.AccessToken,
			AccessTokenExpiresAt:  utils.FormatTimeRFC3339(&session.AccessTokenExpiredAt),
			RefreshToken:          session.RefreshToken,
			RefreshTokenExpiresAt: utils.FormatTimeRFC3339(&session.RefreshTokenExpiredAt),
			TokenType:             "Bearer",
		}
		return c.JSON(http.StatusOK, res)
	}
}

func (s *Service) handleLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctx := c.Request().Context()
		//requester := GetAuthUserFromCtx(ctx)
		//
		//err := s.authUsecase.DeleteSessionByID(c.Request().Context(), requester.SessionID)
		//switch err {
		//case nil:
		//	break
		////case usecase.ErrNotFound:
		////	return ErrNotFound
		//default:
		//	logrus.Error(err)
		//	return httpValidationOrInternalErr(err)
		//}

		return c.NoContent(http.StatusNoContent)
	}
}

// GetAuthUserFromCtx ..
func GetAuthUserFromCtx(ctx context.Context) *model.User {
	authUser := auth.GetUserFromCtx(ctx)
	if authUser == nil {
		return nil
	}
	user := &model.User{
		ID:        authUser.ID,
		SessionID: authUser.SessionID,
	}
	return user
}
