package controller

import (
	"context"
	"kitchen-queue/internal/domain"
)

type KitchenQueueUC interface {
	SendToQueue(ctx context.Context, order domain.Order)
}
