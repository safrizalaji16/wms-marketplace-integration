export interface ApiResponse<T> {
  message: string;
  data: T;
}

export interface ApiErrorPayload {
  message?: string;
  data?: unknown;
}

export interface AuthLoginPayload {
  email: string;
  password: string;
}

export interface AuthLoginData {
  token: string;
}

export interface LogisticChannel {
  id: string;
  name: string;
}

export interface OrderListItemApi {
  id: string;
  order_sn: string;
  wms_status: string;
  marketplace_status: string;
  shipping_status: string;
  tracking_number: string;
  updated_at: string;
}

export interface OrderListDataApi {
  orders: OrderListItemApi[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface OrderItemApi {
  id: string;
  order_id: string;
  sku: string;
  quantity: number;
  price: number;
}

export interface OrderDetailApi {
  id: string;
  order_sn: string;
  shop_id: string;
  wms_status: string;
  marketplace_status: string;
  shipping_status: string;
  tracking_number: string;
  total_amount: number;
  items: OrderItemApi[];
}
