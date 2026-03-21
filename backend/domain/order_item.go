package domain

import "context"

type OrderItem struct {
	ID       string  `bun:"id,pk,type:uuid,default:gen_random_uuid()" db:"id"`
	OrderID  string  `bun:"order_id,notnull" db:"order_id"`
	SKU      string  `bun:"sku,notnull" db:"sku"`
	Quantity int     `bun:"quantity,notnull" db:"quantity"`
	Price    float64 `bun:"price,notnull" db:"price"`
}

type OrderItemRepository interface {
	GetByOrderID(ctx context.Context, orderID string) ([]OrderItem, error)
	Create(ctx context.Context, orderItem *OrderItem) (OrderItem, error)
}
