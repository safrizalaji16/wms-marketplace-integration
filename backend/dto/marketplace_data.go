package dto

type MarketplaceAuthorizeData struct {
	Code string `json:"code"`
}

type MarketplaceTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type MarketplaceTokenRequest struct {
	Code string `json:"code"`
}

type MarketplaceRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ShipOrderRequest struct {
	ChannelID string `json:"channel_id"`
}

type MarketplaceLogisticChannelData struct {
	ChannelID   string `json:"id"`
	ChannelName string `json:"name"`
}

type MarketplaceShipData struct {
	OrderSN        string `json:"order_sn"`
	ShippingStatus string `json:"shipping_status"`
	TrackingNo     string `json:"tracking_no"`
}

type MarketplaceOrderData struct {
	OrderSN        string                     `json:"order_sn" bson:"order_sn"`
	ShopID         string                     `json:"shop_id" bson:"shop_id"`
	Status         string                     `json:"status" bson:"status"`
	ShippingStatus string                     `json:"shipping_status" bson:"shipping_status"`
	TrackingNumber string                     `json:"tracking_number" bson:"tracking_number"`
	TotalAmount    float64                    `json:"total_amount" bson:"total_amount"`
	CreatedAt      string                     `json:"created_at" bson:"created_at"`
	Items          []MarketplaceOrderItemData `json:"items" bson:"items"`
}

type MarketplaceOrderItemData struct {
	SKU      string  `json:"sku" bson:"sku"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
}
