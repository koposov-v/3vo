package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"kitchen-queue/internal/controller/api"
	"kitchen-queue/internal/init/config"
	"kitchen-queue/internal/usecase"
	orderv1 "kitchen-queue/pkg/order/v1"
	v1 "kitchen-queue/pkg/v1"
	"net"
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

	// Инициализация клиента OrderCore
	conn, err := grpc.NewClient(cfg.OrderCoreServiceAddr)
	if err != nil {
		return fmt.Errorf("dial order service: %w", err)
	}
	orderClient := orderv1.NewOrderServiceClient(conn)

	// Инициализация usecase
	kitchenQueueUC := usecase.NewKitchenQueueUseCase(orderClient)

	grpcServer := grpc.NewServer()
	v1.RegisterKitchenServiceServer(
		grpcServer,
		api.NewKitchenServer(
			kitchenQueueUC,
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

	kitchenQueueUC.CloseChannel()

	_, shutdown := context.WithTimeout(context.Background(), cfg.GracefulShutdownDelay)
	defer shutdown()

	grpcServer.GracefulStop()

	return nil
}
