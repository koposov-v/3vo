package domain

import (
	v1 "order-core/pkg/v1"
	"time"
)

type Order struct {
	ID         string
	UserID     string
	Status     int
	Reason     string
	TotalPrice uint32
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Items   []OrderItem
	Comment *string
}

func (o *Order) SetTimestamps() {
	now := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	o.UpdatedAt = now
}

type OrderPatch struct {
	Status     int
	Items      []OrderItem
	TotalPrice uint32
	Comment    *string
}

func (o *Order) Patch(patch OrderPatch) {
	o.Items = patch.Items
	o.Comment = patch.Comment
	o.Status = patch.Status
	o.SetTimestamps()
}

func (o *Order) Cancel() {
	o.Status = int(v1.OrderStatus_ORDER_STATUS_CANCELLED)
	o.SetTimestamps()
}

type OrderItem struct {
	ID       string
	Name     string
	Quantity uint32
	Price    uint32
}
