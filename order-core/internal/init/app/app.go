package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"order-core/internal/controller/api"
	"order-core/internal/init/config"
	"order-core/internal/repository/memory"
	"order-core/internal/usecase"
	v1 "order-core/pkg/v1"
	"os"
	"os/signal"
	"syscall"
)

func Run() error {
	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}

	logger := logrus.New()

	// Инициализация usecase
	orderRepo := memory.NewOrderRepository()
	orderUC := usecase.NewOrderUsecase(orderRepo, logger)

	grpcServer := grpc.NewServer()
	v1.RegisterOrderServiceServer(
		grpcServer,
		api.NewOrderServer(
			orderUC,
			logger,
		),
	)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Grpc.Host, cfg.Grpc.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		logger.Infof("gRPC server listening on %s:%s", cfg.Grpc.Host, cfg.Grpc.Port)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Errorf("gRPC server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	_, shutdown := context.WithTimeout(context.Background(), cfg.GracefulShutdownDelay)
	defer shutdown()

	grpcServer.GracefulStop()

	return nil
}
