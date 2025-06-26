package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"order-core/internal/domain"
	"order-core/internal/repository/memory"
	v1 "order-core/pkg/v1"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	UpdateOrder(ctx context.Context, order *domain.Order) error
	GetOrder(ctx context.Context, id string) (*domain.Order, error)
}

type OrderUsecase struct {
	repo   OrderRepository
	logger *logrus.Logger
}

func NewOrderUsecase(repo *memory.OrderRepository, logger *logrus.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, order *domain.Order) error {
	order.SetTimestamps()

	if err := u.repo.CreateOrder(ctx, order); err != nil {
		u.logger.Errorf("Failed to create order: %v", err)
		return err
	}

	u.SendToKitchen(ctx, order)

	order.Status = int(v1.OrderStatus_ORDER_STATUS_CREATED)
	return nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	return u.repo.GetOrder(ctx, orderID)
}

func (u *OrderUsecase) UpdateOrder(ctx context.Context, order *domain.Order) error {
	if err := u.repo.UpdateOrder(ctx, order); err != nil {
		u.logger.Errorf("Failed to update order: %v", err)
		return err
	}

	return nil
}

func (u *OrderUsecase) CancelOrder(ctx context.Context, order *domain.Order) error {
	order.Cancel()

	if err := u.repo.UpdateOrder(ctx, order); err != nil {
		u.logger.Errorf("Failed to cancel order: %v", err)
		return err
	}

	return nil
}

func (u *OrderUsecase) SendToKitchen(ctx context.Context, order *domain.Order) {
	//kitchenReq := &v1.KitchenOrderRequest{OrderId: order.ID, Items: req.Items}
	//_, err := u.kitchenClient.SendToKitchen(ctx, kitchenReq)
	//if err != nil {
	//	u.logger.Errorf("Failed to send order to kitchen: %v", err)
	//	return err
	//}
}
