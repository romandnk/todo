package app

import (
	"context"
	"github.com/romandnk/todo/config"
	storage "github.com/romandnk/todo/internal/repo"
	httpserver "github.com/romandnk/todo/internal/server/http"
	v1 "github.com/romandnk/todo/internal/server/http/v1"
	"github.com/romandnk/todo/internal/service"
	zaplogger "github.com/romandnk/todo/pkg/logger/zap"
	postgres "github.com/romandnk/todo/pkg/storage"
	"go.uber.org/zap"
	"log"
	"net"
	"os/signal"
	"strconv"
	"syscall"
)

//	@title			TODO App Swagger
//	@version		1.0
//	@description	Swagger API for Golang Project TODO.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API [Roman] Support
//	@license.name	romandnk
//	@license.url	https://github.com/romandnk/todo

// @BasePath	/api/v1/

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	// initializing config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
	}

	// initializing zap logger
	logger, err := zaplogger.NewLogger(cfg.ZapLogger)
	if err != nil {
		log.Fatalf("error initializing zap logger: %s", err.Error())
	}

	logger.Info("using zap logger")

	// initializing connection to postgres db
	db, err := postgres.NewStorage(ctx, cfg.Postgres)
	if err != nil {
		db.Close()
		logger.Fatal("error initializing postgres db", zap.Error(err))
	}
	defer db.Close()

	logger.Info("using postgres repo",
		zap.String("host", cfg.Postgres.Host),
		zap.Int("port", cfg.Postgres.Port),
	)

	// initializing repository
	repo := storage.NewRepository(db)

	// initializing service dependencies
	dep := service.Dependencies{
		Repo:   repo,
		Logger: logger,
	}

	// initializing services
	services := service.NewServices(dep)

	// initializing middlewares
	mw := v1.NewMiddlewares(logger)

	// initializing http handler
	handler := v1.NewHandler(services, logger, mw)

	// initializing http server
	srv := httpserver.NewServer(cfg.HTTPServer, handler.InitRoutes())

	logger.Info("starting http server...",
		"address", net.JoinHostPort(cfg.HTTPServer.Host, strconv.Itoa(cfg.HTTPServer.Port)))
	srv.Start()

	select {
	case <-ctx.Done():
		logger.Info("stopping http server...")

		err = srv.Stop(ctx)
		if err != nil {
			logger.Error("error stopping http server", zap.Error(err))
		}

		logger.Info("http server is stopped")
	case err = <-srv.Notify():
		logger.Error("error starting http server", zap.Error(err))
		cancel()
	}
}
