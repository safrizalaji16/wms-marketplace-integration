package dto

type OrderItemData struct {
	ID       string  `json:"id"`
	OrderID  string  `json:"order_id"`
	SKU      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
