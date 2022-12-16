package httpsvc

import (
	"github.com/irvankadhafi/user-balance-transfer-service/auth"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Service http service
type Service struct {
	echo           *echo.Echo
	authUsecase    model.AuthUsecase
	httpMiddleware *auth.AuthenticationMiddleware
}

// RouteService add dependencies and use group for routing
func RouteService(
	echo *echo.Echo,
	authUsecase model.AuthUsecase,
	authMiddleware *auth.AuthenticationMiddleware,
) {
	srv := &Service{
		echo:           echo,
		authUsecase:    authUsecase,
		httpMiddleware: authMiddleware,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	// auth
	s.echo.POST("/auth/login/", s.handleLoginByUsernamePassword())
	s.echo.GET("/irvan/", s.handleIrvan(), s.httpMiddleware.MustAuthenticateAccessToken())

}

func (s *Service) handleIrvan() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Warn("MASUK HANDLE IRVAN")
		return c.JSON(http.StatusOK, "HAI IRVAN")
	}
}
