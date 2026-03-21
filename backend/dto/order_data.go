package dto

type OrderData struct {
	ID                string `json:"id"`
	OrderSN           string `json:"order_sn"`
	WMSStatus         string `json:"wms_status"`
	MarketplaceStatus string `json:"marketplace_status"`
	ShippingStatus    string `json:"shipping_status"`
	TrackingNumber    string `json:"tracking_number"`
	UpdatedAt         string `json:"updated_at"`
}

type OrderFilter struct {
	Search            string `json:"search"`
	WMSStatus         string `json:"wms_status"`
	MarketplaceStatus string `json:"marketplace_status"`
	ShippingStatus    string `json:"shipping_status"`
	SortBy            string `json:"sort_by"`
	SortDir           string `json:"sort_dir"`
	Page              int    `json:"page"`
	Limit             int    `json:"limit"`
	Offset            int    `json:"-"`
}

type OrderListData struct {
	Orders     []OrderData `json:"orders"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type OrderDetailData struct {
	ID                string          `json:"id"`
	OrderSN           string          `json:"order_sn"`
	ShopID            string          `json:"shop_id"`
	WMSStatus         string          `json:"wms_status"`
	MarketplaceStatus string          `json:"marketplace_status"`
	ShippingStatus    string          `json:"shipping_status"`
	TrackingNumber    string          `json:"tracking_number"`
	TotalAmount       float64         `json:"total_amount"`
	Items             []OrderItemData `json:"items"`
}

type OrderStatusWebhookRequest struct {
	OrderSN string `json:"order_sn"`
	Status  string `json:"status"`
}

type ShippingStatusWebhookRequest struct {
	OrderSN string `json:"order_sn"`
	Status  string `json:"status"`
}
