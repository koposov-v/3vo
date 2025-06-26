package app

import (
	"api-gateway/internal/adapter/srv"
	"api-gateway/internal/controller/api"
	"api-gateway/internal/init/config"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderv1 "kitchen-queue/pkg/order/v1"
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

	// Инициализация клиента OrderCore
	conn, err := grpc.NewClient(
		cfg.OrderCoreServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("dial order service: %w", err)
	}
	orderClient := orderv1.NewOrderServiceClient(conn)

	gateway := api.NewGateway(
		orderClient,
		cfg.AuthServiceURL,
		logger,
	)

	//main
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Роуты
	e.POST("/login", gateway.Login)

	orderGroup := e.Group("/order", gateway.AuthMiddleware())
	orderGroup.GET("/:order_id", gateway.GetOrder)
	orderGroup.POST("", gateway.CreateOrder)
	orderGroup.PATCH("", gateway.UpdateOrder)
	orderGroup.DELETE("/:order_id", gateway.CancelOrder)

	logger.Infof("Start server on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
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
