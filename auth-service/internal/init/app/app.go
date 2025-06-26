package app

import (
	"authjwt/internal/adapter/srv"
	v1 "authjwt/internal/controller/api/v1"
	"authjwt/internal/init/config"
	"authjwt/internal/repository/memory"
	"authjwt/internal/usecase"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
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

	userUseCase := usecase.NewUserUseCase(
		memory.NewInMemoryUserRepo(),
		cfg.JWTSecret,
	)

	//main
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрация маршрутов
	userRoutes := v1.NewUserRoutes(
		userUseCase,
		v10,
	)

	groupV1 := e.Group("/api/v1")
	userRoutes.Register(groupV1)
	server := srv.NewServer(cfg, e)

	for _, r := range e.Routes() {
		logger.Infof(fmt.Sprintf("%s %s", r.Method, r.Path))
	}

	go func() {
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), cfg.GracefulShutdownDelay)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Errorf("error occurred while shutting down http server: %s\n", err.Error())
	}

	return nil
}
