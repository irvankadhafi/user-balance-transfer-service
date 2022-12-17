package console

import (
	"fmt"
	"github.com/irvankadhafi/user-balance-transfer-service/auth"
	"github.com/irvankadhafi/user-balance-transfer-service/cacher"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/config"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/db"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/delivery/httpsvc"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/helper"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/repository"
	"github.com/irvankadhafi/user-balance-transfer-service/internal/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	// Initiate all connection like db, redis, etc
	db.InitializePostgresConn()
	authenticationCacher := cacher.NewCacheManager()
	generalCacher := cacher.NewCacheManager()
	pgDB, err := db.PostgreSQL.DB()
	continueOrFatal(err)
	defer helper.WrapCloser(pgDB.Close)

	redisOpts := &RedisConnectionPoolOptions{
		DialTimeout:     config.RedisDialTimeout(),
		ReadTimeout:     config.RedisReadTimeout(),
		WriteTimeout:    config.RedisWriteTimeout(),
		IdleCount:       config.RedisMaxIdleConn(),
		PoolSize:        config.RedisMaxActiveConn(),
		IdleTimeout:     240 * time.Second,
		MaxConnLifetime: 1 * time.Minute,
	}

	authRedisConn, err := NewRedigoRedisConnectionPool(config.RedisAuthCacheHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(authRedisConn.Close)

	authRedisLockConn, err := NewRedigoRedisConnectionPool(config.RedisAuthCacheLockHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(authRedisLockConn.Close)

	authenticationCacher.SetConnectionPool(authRedisConn)
	authenticationCacher.SetLockConnectionPool(authRedisLockConn)
	authenticationCacher.SetDefaultTTL(config.CacheTTL())

	redisConn, err := NewRedigoRedisConnectionPool(config.RedisCacheHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(redisConn.Close)

	redisLockConn, err := NewRedigoRedisConnectionPool(config.RedisLockHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(redisLockConn.Close)

	generalCacher.SetConnectionPool(redisConn)
	generalCacher.SetLockConnectionPool(redisLockConn)
	generalCacher.SetDefaultTTL(config.CacheTTL())

	userRepo := repository.NewUserRepository(db.PostgreSQL, generalCacher)
	userUsecase := usecase.NewUserUsecase(userRepo)

	sessionRepo := repository.NewSessionRepository(db.PostgreSQL, authenticationCacher, userRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, userUsecase)
	userAuther := usecase.NewUserAutherAdapter(authUsecase)

	gormTransationer := repository.NewGormTransactioner(db.PostgreSQL)
	userBalanceRepo := repository.NewUserBalanceRepository(db.PostgreSQL)
	userBalanceHistoryRepo := repository.NewUserBalanceHistoryRepository(db.PostgreSQL)
	userBalanceUsecase := usecase.NewUserBalanceUsecase(userRepo, userBalanceRepo, userBalanceHistoryRepo, gormTransationer, sessionRepo)

	bankBalanceRepo := repository.NewBankBalanceRepository(db.PostgreSQL)
	bankBalanceHistoryRepo := repository.NewBankBalanceHistoryRepository(db.PostgreSQL)
	bankBalanceUsecase := usecase.NewBankBalanceUsecase(bankBalanceRepo, bankBalanceHistoryRepo, gormTransationer, sessionRepo)

	httpServer := echo.New()
	httpMiddleware := auth.NewAuthenticationMiddleware(authenticationCacher, userAuther)

	httpServer.Pre(middleware.AddTrailingSlash())
	httpServer.Use(middleware.Logger())
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.CORS())

	httpsvc.RouteService(httpServer, authUsecase, userBalanceUsecase, bankBalanceUsecase, httpMiddleware)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	quitCh := make(chan bool, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for {
			select {
			case <-sigCh:
				gracefulShutdown(httpServer)
				quitCh <- true
			case e := <-errCh:
				log.Error(e)
				gracefulShutdown(httpServer)
				quitCh <- true
			}
		}
	}()

	go func() {
		// Start HTTP server
		if err := httpServer.Start(fmt.Sprintf(":%s", config.HTTPPort())); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	<-quitCh
	log.Info("exiting")
}
