package dto

type CreateOrderRequest struct {
	UserID  string      `json:"user_id"`
	Items   []OrderItem `json:"items"`
	Comment *string     `json:"comment,omitempty"`
}

type OrderItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity uint32 `json:"quantity"`
	Price    uint32 `json:"price"`
}

type UpdateOrderRequest struct {
	OrderID string      `json:"order_id"`
	Items   []OrderItem `json:"items,omitempty"`
	Comment *string     `json:"comment,omitempty"`
	Status  string      `json:"status"`
}

type CancelOrderBody struct {
	Reason string `json:"reason"`
}
