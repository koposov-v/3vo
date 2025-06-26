package domain

type Order struct {
	ID      string
	Items   []OrderItem
	Comment *string
}

type OrderItem struct {
	ID       string
	Name     string
	Quantity uint32
	Price    uint32
}
