package domain

import (
	"backend/dto"
	"context"
	"database/sql"
	"encoding/json"
)

type Order struct {
	ID                    string          `bun:"id,pk,type:uuid,default:gen_random_uuid()" db:"id"`
	OrderSN               string          `bun:"order_sn,notnull" db:"order_sn"`
	ShopID                string          `bun:"shop_id,notnull" db:"shop_id"`
	MarketplaceStatus     string          `bun:"marketplace_status" db:"marketplace_status"`
	ShippingStatus        string          `bun:"shipping_status" db:"shipping_status"`
	WMSStatus             string          `bun:"wms_status" db:"wms_status"`
	TrackingNumber        string          `bun:"tracking_number" db:"tracking_number"`
	TotalAmount           float64         `bun:"total_amount" db:"total_amount"`
	RawMarketplacePayload json.RawMessage `bun:"raw_marketplace_payload,type:jsonb" db:"raw_marketplace_payload"`
	CreatedAt             sql.NullTime    `bun:"created_at,nullzero,notnull,default:current_timestamp" db:"created_at"`
	UpdatedAt             sql.NullTime    `bun:"updated_at,nullzero,notnull,default:current_timestamp" db:"updated_at"`
	Items                 []OrderItem     `bun:"rel:has-many,join:id=order_id" json:"-"`
}

type OrderRepository interface {
	Count(ctx context.Context) (int, error)
	GetAll(ctx context.Context, filter dto.OrderFilter) ([]Order, int, error)
	GetByOrderSN(ctx context.Context, orderSN string) (Order, error)
	Create(ctx context.Context, order *Order) (Order, error)
	Update(ctx context.Context, order *Order) error
}

type OrderService interface {
	GetOrders(ctx context.Context, filter dto.OrderFilter) (dto.OrderListData, error)
	GetByOrderSN(ctx context.Context, orderSN string) (dto.OrderDetailData, error)
	Pick(ctx context.Context, orderSN string) error
	Pack(ctx context.Context, orderSN string) error
	Ship(ctx context.Context, orderSN, channelID string) error
	UpdateMarketplaceStatus(ctx context.Context, orderSN, status string) error
	UpdateShippingStatus(ctx context.Context, orderSN, status string) error
}
