package httpsvc

import (
	"github.com/irvankadhafi/user-balance-transfer-service/auth"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/model"
	"github.com/irvankadhafi/user-balance-transfer-service/utils"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Service http service
type Service struct {
	echo               *echo.Echo
	authUsecase        model.AuthUsecase
	userBalanceUsecase model.UserBalanceUsecase
	bankBalanceUsecase model.BankBalanceUsecase
	httpMiddleware     *auth.AuthenticationMiddleware
}

// RouteService add dependencies and use group for routing
func RouteService(
	echo *echo.Echo,
	authUsecase model.AuthUsecase,
	userBalanceUsecase model.UserBalanceUsecase,
	bankBalanceUsecase model.BankBalanceUsecase,
	authMiddleware *auth.AuthenticationMiddleware,
) {
	srv := &Service{
		echo:               echo,
		authUsecase:        authUsecase,
		userBalanceUsecase: userBalanceUsecase,
		bankBalanceUsecase: bankBalanceUsecase,
		httpMiddleware:     authMiddleware,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	// auth
	s.echo.POST("/auth/login/", s.handleLoginByEmailPassword())

	s.echo.GET("/user-balance/", s.handleGetUserBalance(), s.httpMiddleware.MustAuthenticateAccessToken())
	s.echo.POST("/user-balance/add/", s.handleAddUserBalance(), s.httpMiddleware.MustAuthenticateAccessToken())
	s.echo.POST("/user-balance/transfer/", s.handleUserBalanceTransfer(), s.httpMiddleware.MustAuthenticateAccessToken())

	s.echo.POST("/bank-balance/create/", s.handleCreateBankBalance(), s.httpMiddleware.MustAuthenticateAccessToken())
	s.echo.POST("/bank-balance/add/", s.handleAddBankBalance(), s.httpMiddleware.MustAuthenticateAccessToken())
	s.echo.POST("/bank-balance/transfer/", s.handleAddBankBalance(), s.httpMiddleware.MustAuthenticateAccessToken())

}

func (s *Service) handleAddBankBalance() echo.HandlerFunc {
	type request struct {
		Code    string `json:"code"`
		Balance string `json:"balance"`
		Author  string `json:"author"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := GetAuthUserFromCtx(ctx)

		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		err := s.bankBalanceUsecase.AddBankBalance(ctx, model.AddBankBalanceInput{
			Code:      req.Code,
			Balance:   utils.StringToInt64(req.Balance),
			SessionID: user.SessionID,
			Author:    req.Author,
		})
		switch err {
		case nil:
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusOK, "ok")
	}
}

func (s *Service) handleCreateBankBalance() echo.HandlerFunc {
	type request struct {
		Code string `json:"code"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := GetAuthUserFromCtx(ctx)

		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		err := s.bankBalanceUsecase.CreateBankAccount(ctx, user.SessionID, req.Code)
		switch err {
		case nil:
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusOK, "ok")
	}
}

func (s *Service) handleGetUserBalance() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := GetAuthUserFromCtx(ctx)

		balance, err := s.userBalanceUsecase.GetCurrentUserBalanceByUserID(ctx, user.ID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, balance)
	}
}

func (s *Service) handleUserBalanceTransfer() echo.HandlerFunc {
	type request struct {
		ToUserID int    `json:"to_user_id"`
		Balance  string `json:"balance"`
		Author   string `json:"author"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := GetAuthUserFromCtx(ctx)

		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}
		err := s.userBalanceUsecase.TransferUserBalance(ctx, model.TransferUserBalanceInput{
			FromUserID: user.ID,
			ToUserID:   req.ToUserID,
			Balance:    utils.StringToInt64(req.Balance),
			SessionID:  user.SessionID,
			Author:     req.Author,
		})
		switch err {
		case nil:
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusOK, "ok")
	}
}

func (s *Service) handleAddUserBalance() echo.HandlerFunc {
	type request struct {
		Balance string `json:"balance"`
		Author  string `json:"author"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := GetAuthUserFromCtx(ctx)

		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		err := s.userBalanceUsecase.AddUserBalance(ctx, model.AddUserBalanceInput{
			UserID:    user.ID,
			Balance:   utils.StringToInt64(req.Balance),
			SessionID: user.SessionID,
			Author:    req.Author,
		})
		switch err {
		case nil:
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusOK, "ok")
	}
}
