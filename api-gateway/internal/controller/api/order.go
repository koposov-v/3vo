package api

import (
	"api-gateway/internal/controller/dto"
	orderv1 "api-gateway/pkg/order/v1"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (g *Gateway) GetOrder(c echo.Context) error {
	g.logger.Info("Заход на маршрут GET /order")

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order_id is required")
	}

	grpcReq := &orderv1.GetOrderRequest{
		OrderId: orderID,
	}

	ctx := c.Request().Context()
	resp, err := g.orderClient.GetOrder(ctx, grpcReq)
	if err != nil {
		g.logger.Error("ошибка при вызове gRPC GetOrder", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (g *Gateway) CreateOrder(c echo.Context) error {
	g.logger.Info("Заход на маршрут POST /order")

	var req dto.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		g.logger.Error("не удалось распарсить тело запроса", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	grpcReq := toGRPCCreateOrderRequest(req)

	ctx := c.Request().Context()
	resp, err := g.orderClient.CreateOrder(ctx, grpcReq)
	if err != nil {
		g.logger.Error("ошибка вызова gRPC CreateOrder", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create order")
	}

	return c.JSON(http.StatusOK, resp)
}

func (g *Gateway) UpdateOrder(c echo.Context) error {
	g.logger.Info("Заход на маршрут PUT /order")

	var req dto.UpdateOrderRequest
	if err := c.Bind(&req); err != nil {
		g.logger.Error("не удалось распарсить тело запроса", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	grpcReq, err := toGRPCUpdateOrderRequest(req)
	if err != nil {
		g.logger.Error("ошибка маппинга в gRPC UpdateOrderRequest", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	resp, err := g.orderClient.UpdateOrder(ctx, grpcReq)
	if err != nil {
		g.logger.Error("ошибка вызова gRPC UpdateOrder", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update order")
	}

	return c.JSON(http.StatusOK, resp)
}

func (g *Gateway) CancelOrder(c echo.Context) error {
	g.logger.Info("Заход на маршрут DELETE /order/:order_id")

	orderID := c.Param("order_id")
	if orderID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order_id is required")
	}

	var body dto.CancelOrderBody
	if err := c.Bind(&body); err != nil {
		g.logger.Error("не удалось распарсить reason из тела", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	grpcReq := &orderv1.CancelOrderRequest{
		OrderId: orderID,
		Reason:  body.Reason,
	}

	ctx := c.Request().Context()
	_, err := g.orderClient.CancelOrder(ctx, grpcReq)
	if err != nil {
		g.logger.Error("ошибка при вызове gRPC CancelOrder", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to cancel order")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order cancelled"})
}

func toGRPCCreateOrderRequest(in dto.CreateOrderRequest) *orderv1.CreateOrderRequest {
	items := make([]*orderv1.OrderItem, 0, len(in.Items))
	for _, i := range in.Items {
		items = append(items, &orderv1.OrderItem{
			Id:       i.ID,
			Name:     i.Name,
			Quantity: i.Quantity,
			Price:    i.Price,
		})
	}

	req := &orderv1.CreateOrderRequest{
		UserId: in.UserID,
		Items:  items,
	}
	if in.Comment != nil {
		req.Comment = *in.Comment
	}
	return req
}

func toGRPCUpdateOrderRequest(req dto.UpdateOrderRequest) (*orderv1.UpdateOrderRequest, error) {
	status, ok := orderv1.OrderStatus_value[req.Status]
	if !ok {
		return nil, fmt.Errorf("unknown status: %s", req.Status)
	}

	items := make([]*orderv1.OrderItem, 0, len(req.Items))
	for _, i := range req.Items {
		items = append(items, &orderv1.OrderItem{
			Id:       i.ID,
			Name:     i.Name,
			Quantity: i.Quantity,
			Price:    i.Price,
		})
	}

	return &orderv1.UpdateOrderRequest{
		OrderId: req.OrderID,
		Items:   items,
		Comment: req.Comment,
		Status:  orderv1.OrderStatus(status),
	}, nil
}
