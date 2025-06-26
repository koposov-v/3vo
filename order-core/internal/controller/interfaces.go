package controller

import (
	"context"
	"order-core/internal/domain"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrder(ctx context.Context, orderID string) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	CancelOrder(ctx context.Context, order *domain.Order) error
}
