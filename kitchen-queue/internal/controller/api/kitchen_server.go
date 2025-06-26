package api

import (
	"context"
	"github.com/sirupsen/logrus"
	"kitchen-queue/internal/controller"
	"kitchen-queue/internal/domain"
	v1 "kitchen-queue/pkg/v1"
)

type KitchenServer struct {
	v1.UnimplementedKitchenServiceServer
	kitchenQueueUC controller.KitchenQueueUC
	orderClient    any
	logger         *logrus.Logger
}

func NewKitchenServer(kitchenQueueUC controller.KitchenQueueUC, logger *logrus.Logger) *KitchenServer {
	return &KitchenServer{
		kitchenQueueUC: kitchenQueueUC,
		logger:         logger,
	}
}

func (s *KitchenServer) SendToKitchen(ctx context.Context, req *v1.SendToKitchenRequest) (*v1.KitchenStatusResponse, error) {
	items := fromItems(req.Items)
	orderKitchen := domain.Order{
		ID:      req.OrderId,
		Items:   items,
		Comment: req.Comment,
	}

	s.kitchenQueueUC.SendToQueue(ctx, orderKitchen)

	return &v1.KitchenStatusResponse{}, nil
}

func fromItems(req []*v1.OrderItem) []domain.OrderItem {
	items := make([]domain.OrderItem, len(req))
	for i, item := range req {
		items[i] = domain.OrderItem{
			ID:       item.Id,
			Name:     item.Name,
			Quantity: item.Quantity,
		}
	}
	return items
}
