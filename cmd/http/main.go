package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/config"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/http/handlers"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/http/middlewares"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/dummy"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/pgsql"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/redis"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/smsgateway"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/AshkanAbd/arvancloud_sms_gateway/docs"
	smssrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	usersrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/services"
	pkgCfg "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/config"
	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
	pkgMetrics "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/metrics"
	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
	pkgRedis "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/redis"
)

var Config config.AppConfig

func init() {
	pkgLog.SetLogLevel("Trace")
	Config = *pkgCfg.Load[config.AppConfig]()
	pkgLog.SetLogLevel(Config.LogLevel)

	pkgLog.Info("Config is %#v", Config)
}

func closer(c io.Closer) {
	if err := c.Close(); err != nil {
		pkgLog.Error(err, "close error")
	}
}

// @title			SMS Gateway API
// @version		1.0
// @description	SMS Gateway API
// @host			localhost:8000
// @BasePath		/
func main() {
	appCtx, cancel := context.WithCancel(context.Background())

	pgsqlConn, err := pkgPgSql.NewConnector(Config.PgSQLConfig)
	if err != nil {
		pkgLog.Error(err, "failed to connect to pgsql")
		return
	}
	defer closer(pgsqlConn)

	if err := pgsqlConn.Migrate(); err != nil {
		pkgLog.Error(err, "failed to migrate pgsql")
		return
	}

	redisConn := pkgRedis.NewConnector(Config.RedisConfig)
	defer closer(redisConn)

	pgsqlRepo, err := pgsql.NewRepository(pgsqlConn)
	if err != nil {
		pkgLog.Error(err, "failed to create pgsql repository")
		return
	}

	redisRepo := redis.NewRepository(Config.RedisRepoConfig, redisConn)

	smsSender := dummy.NewSmsSender()

	userService := usersrv.NewUserService(pgsqlRepo)
	smsService := smssrv.NewSmsService(Config.SmsServiceConfig, pgsqlRepo, smsSender, redisRepo)

	gateway := smsgateway.NewSmsGateway(Config.SmsGatewayConfig, userService, smsService)

	pkgMetrics.RegisterMetrics()

	httpHandler := handlers.NewHttpHandler(gateway)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(middlewares.Logger())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Post("/user", httpHandler.CreateUser)
	app.Get("/user/:id", httpHandler.GetUser)
	app.Get("/user/:id/sms", httpHandler.GetUserMessages)
	app.Post("/user/:id/balance", httpHandler.IncreaseUserBalance)
	app.Post("/user/:id/sms/single", httpHandler.SendSingleMessage)
	app.Post("/user/:id/sms/bulk", httpHandler.SendBulkMessage)

	enqueueWorkerErrCh := gateway.StartEnqueueWorker(appCtx)
	sendWorkerErrCh := gateway.StartSendWorkers(appCtx)
	httpErrCh := make(chan error, 1)

	app.Get("/healthz", handlers.HealthCheck([]<-chan error{
		enqueueWorkerErrCh,
		sendWorkerErrCh,
		httpErrCh,
	}))

	app.Get("/metrics", handlers.Metrics())

	go func() {
		pkgLog.Debug("Starting http server on %s", Config.HttpConfig.Address)
		if err := app.Listen(Config.HttpConfig.Address); err != nil {
			httpErrCh <- err
			pkgLog.Error(err, "Error on starting http server")
		}
	}()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	stop := <-signalForExit
	pkgLog.Info("Stop signal received: %v", stop)
	pkgLog.Info("Shutting down...")
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		pkgLog.Error(err, "Error on shutdown http server after 5 seconds")
	}
	cancel()
}
