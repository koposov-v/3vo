package memory

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"order-core/internal/domain"
	"sync"
)

type OrderRepository struct {
	orders map[string]*domain.Order
	mutex  sync.RWMutex
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *OrderRepository) CreateOrder(_ context.Context, order *domain.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	order.ID = uuid.NewString()

	r.orders[order.ID] = order

	return nil
}

func (r *OrderRepository) GetOrder(_ context.Context, orderID string) (*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}

	return order, nil
}

func (r *OrderRepository) UpdateOrder(_ context.Context, order *domain.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return fmt.Errorf("order with ID %s not found", order.ID)
	}
	r.orders[order.ID] = order
	return nil
}
