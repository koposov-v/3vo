package usecase

import (
	"context"
	"fmt"
	"kitchen-queue/internal/domain"
	orderv1 "kitchen-queue/pkg/order/v1"
	"sync"
	"time"
)

var statuses = []orderv1.OrderStatus{
	orderv1.OrderStatus_ORDER_STATUS_QUEUED,
	orderv1.OrderStatus_ORDER_STATUS_PREPARING,
	orderv1.OrderStatus_ORDER_STATUS_READY,
}

type KitchenQueueUseCase struct {
	wg          *sync.WaitGroup
	queue       chan domain.Order
	orderClient orderv1.OrderServiceClient
}

func NewKitchenQueueUseCase(orderClient orderv1.OrderServiceClient) *KitchenQueueUseCase {
	return &KitchenQueueUseCase{
		queue:       make(chan domain.Order),
		orderClient: orderClient,
		wg:          &sync.WaitGroup{},
	}
}

func (uc *KitchenQueueUseCase) SendToQueue(ctx context.Context, order domain.Order) {
	uc.queue <- order
}

func (uc *KitchenQueueUseCase) CloseChannel() {
	uc.wg.Wait()
	close(uc.queue)
}

func (uc *KitchenQueueUseCase) worker() {
	for order := range uc.queue {
		uc.wg.Add(1)
		go func() {
			defer uc.wg.Done()
			for _, st := range statuses {
				//Имитируем работу, проходился по каждому статусу и выполняем работу -> отправляем в OrderCore
				time.Sleep(30 * time.Second)
				_, err := uc.orderClient.UpdateOrder(context.Background(), &orderv1.UpdateOrderRequest{
					Status:  st,
					OrderId: order.ID,
				})

				if err != nil {
					fmt.Printf("failed to update order %s to status %v: %v\n", order.ID, st, err)
				}
			}
		}()
	}
}
