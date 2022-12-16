package auth

import (
	"context"
	"github.com/irvankadhafi/user-balance-transfer-service/cacher"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

// UserAuthenticator to perform user authentication
type UserAuthenticator interface {
	AuthenticateToken(ctx context.Context, accessToken string) (*User, error)
}

// AuthenticationMiddleware middleware for authentication
type AuthenticationMiddleware struct {
	cacheManager      cacher.CacheManager
	userAuthenticator UserAuthenticator
}

// NewAuthenticationMiddleware AuthMiddleware constructor
func NewAuthenticationMiddleware(
	manager cacher.CacheManager,
	userAuthenticator UserAuthenticator,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		cacheManager:      manager,
		userAuthenticator: userAuthenticator,
	}
}

// AuthenticateAccessToken authenticate access token from http `Authorization` header and load a User to context
func (a *AuthenticationMiddleware) AuthenticateAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getAccessToken(c.Request())
			return a.authenticateAccessToken(c, next, token)
		}
	}
}

// MustAuthenticateAccessToken must authenticate access token from http `Authorization` header and load a User to context
// Differ from AuthenticateAccessToken, if no token provided then return Unauthenticated
func (a *AuthenticationMiddleware) MustAuthenticateAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getAccessToken(c.Request())
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "user is unauthenticated"})
			}

			return a.authenticateAccessToken(c, next, token)
		}
	}
}

func getAccessToken(req *http.Request) (accessToken string) {
	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func (a *AuthenticationMiddleware) authenticateAccessToken(c echo.Context, next echo.HandlerFunc, token string) error {
	// only load user to context when token presented
	if token == "" {
		return next(c)
	}
	ctx := c.Request().Context()

	session, err := a.findSessionFromCache(token)
	switch err {
	default:
		logrus.WithField("sessionCacheError", "find session from cache got error").Error(err)
	case nil:
		if session == nil {
			break // fallback
		}
		if session.IsAccessTokenExpired() {
			return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "token expired"})
		}

		ctx := SetUserToCtx(c.Request().Context(), NewUserFromSession(*session))
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}

	userSession, err := a.userAuthenticator.AuthenticateToken(ctx, token)
	// fallback to rpc
	switch status.Code(err) {
	case codes.OK:
		if userSession == nil { // safety check
			return next(c)
		}

		ctx := SetUserToCtx(c.Request().Context(), *userSession)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	case codes.NotFound:
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "token is invalid"})
	case codes.Unauthenticated:
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "token is expired"})
	default:
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": "system error"})
	}
}

func (a *AuthenticationMiddleware) findSessionFromCache(token string) (*model.Session, error) {
	reply, err := a.cacheManager.Get(model.NewSessionTokenCacheKey(token))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if reply == nil {
		return nil, nil
	}

	sess := utils.InterfaceBytesToType[*model.Session](reply)
	return sess, nil
}
