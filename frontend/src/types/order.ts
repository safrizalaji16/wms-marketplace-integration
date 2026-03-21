export type MarketplaceStatus =
  | "processing"
  | "paid"
  | "shipping"
  | "delivered"
  | "cancelled";

export type ShippingStatus =
  | "label_created"
  | "awaiting_pickup"
  | "shipped"
  | "delivered"
  | "cancelled";

export type WmsStatus = "READY_TO_PICK" | "PICKING" | "PACKED" | "SHIPPED";

export type SortDirection = "asc" | "desc";

export interface OrderItem {
  sku: string;
  quantity: number;
  price: number;
}

export interface Order {
  order_sn: string;
  shop_id: string;
  marketplace_status: MarketplaceStatus;
  shipping_status: ShippingStatus;
  wms_status: WmsStatus;
  tracking_number: string | null;
  total_amount: number;
  created_at: string;
  updated_at: string;
  raw_marketplace_payload: Record<string, unknown>;
  items: OrderItem[];
}

export interface OrdersResponse {
  orders: Order[];
  page?: number;
  limit?: number;
  total?: number;
  total_pages?: number;
}
