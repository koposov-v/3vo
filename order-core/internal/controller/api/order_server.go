package api

import (
	"context"
	"github.com/sirupsen/logrus"
	"order-core/internal/controller"
	"order-core/internal/domain"
	v1 "order-core/pkg/v1"
	"time"
)

type OrderServer struct {
	v1.UnimplementedOrderServiceServer
	orderUC controller.OrderUsecase
	logger  *logrus.Logger
}

func NewOrderServer(orderUC controller.OrderUsecase, logger *logrus.Logger) *OrderServer {
	return &OrderServer{
		orderUC: orderUC,
		logger:  logger,
	}
}

func (c *OrderServer) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.OrderResponse, error) {
	order := fromCreateOrderRequest(req)

	err := c.orderUC.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(order), err
}

func (c *OrderServer) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.OrderResponse, error) {
	order, err := c.orderUC.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return toOrderResponse(order), err
}

func (c *OrderServer) UpdateOrder(ctx context.Context, req *v1.UpdateOrderRequest) (*v1.OrderResponse, error) {
	//Todo::можно добавлять разную валидацию еще, но я думаю это и так понятно
	c.logger.Infof("UpdateOrder")
	originalOrder, err := c.orderUC.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	patch := fromUpdateOrderRequest(req)
	originalOrder.Patch(patch)

	if err := c.orderUC.UpdateOrder(ctx, originalOrder); err != nil {
		return nil, err
	}

	return toOrderResponse(originalOrder), err
}

func (c *OrderServer) CancelOrder(ctx context.Context, req *v1.CancelOrderRequest) (*v1.OrderResponse, error) {
	order, err := c.orderUC.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	order.Reason = req.Reason
	if err := c.orderUC.CancelOrder(ctx, order); err != nil {
		return nil, err
	}

	return toOrderResponse(order), err
}

func fromCreateOrderRequest(req *v1.CreateOrderRequest) *domain.Order {
	items, totalPrice := fromItems(req.Items)
	return &domain.Order{
		UserID:     req.UserId,
		Items:      items,
		Comment:    &req.Comment,
		Status:     int(v1.OrderStatus_ORDER_STATUS_CREATED), // TODO::надо бы сделать, domain енамы, но пока лень
		TotalPrice: totalPrice,
	}
}
func fromItems(req []*v1.OrderItem) ([]domain.OrderItem, uint32) {
	var totalPrice uint32
	items := make([]domain.OrderItem, len(req))
	for i, item := range req {
		items[i] = domain.OrderItem{
			ID:       item.Id,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
		totalPrice += item.Price * item.Quantity
	}
	return items, totalPrice
}

func fromUpdateOrderRequest(req *v1.UpdateOrderRequest) domain.OrderPatch {
	items, totalPrice := fromItems(req.Items)
	return domain.OrderPatch{
		Status:     int(req.Status),
		Items:      items,
		Comment:    req.Comment,
		TotalPrice: totalPrice,
	}
}

func toOrderResponse(order *domain.Order) *v1.OrderResponse {
	return &v1.OrderResponse{
		Id:         order.ID,
		UserId:     order.UserID,
		Items:      toItems(order.Items),
		Comment:    order.Comment,
		Status:     v1.OrderStatus(order.Status),
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
		TotalPrice: order.TotalPrice,
	}
}

func toItems(orders []domain.OrderItem) []*v1.OrderItem {
	items := make([]*v1.OrderItem, len(orders))
	for i, orderItem := range orders {
		items[i] = &v1.OrderItem{
			Id:       orderItem.ID,
			Name:     orderItem.Name,
			Quantity: orderItem.Quantity,
			Price:    orderItem.Price,
		}
	}
	return items
}
