package app

import (
	"context"
	"github.com/romandnk/todo/config"
	postgres "github.com/romandnk/todo/pkg/storage"
	"log"
	"os/signal"
	"syscall"
)

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

	logger.Info("using postgres storage",
		zap.String("host", cfg.Postgres.Host),
		zap.Int("port", cfg.Postgres.Port),
	)
}
