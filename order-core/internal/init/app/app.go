package app

import (
	_ "authjwt/docs"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"order-core/internal/init/config"
	"os"
	"os/signal"
	"syscall"
)

func Run() error {
	ctx := context.Background()

	cfg, err := config.Init()
	if err != nil {
		return err
	}

	logger := logrus.New()
	v10 := validator.New()
	_, _, _ = v10, logger, ctx

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), cfg.GracefulShutdownDelay)
	defer shutdown()

	return nil
}
