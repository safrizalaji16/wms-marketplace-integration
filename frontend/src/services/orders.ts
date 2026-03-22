import { apiFetch } from "./api";
import type {
  LogisticChannel,
  OrderDetailApi,
  OrderItemApi,
  OrderListDataApi,
  OrderListItemApi
} from "../types/api";
import type { Order, OrdersResponse, SortDirection, WmsStatus } from "../types/order";

const ordersApiUrl = "/api/orders";

function normalizeOrderListItem(order: OrderListItemApi): Order {
  return {
    order_sn: order.order_sn,
    shop_id: "",
    marketplace_status: order.marketplace_status as Order["marketplace_status"],
    shipping_status: order.shipping_status as Order["shipping_status"],
    wms_status: order.wms_status as WmsStatus,
    tracking_number: order.tracking_number || null,
    total_amount: 0,
    created_at: "",
    updated_at: order.updated_at,
    raw_marketplace_payload: {},
    items: []
  };
}

function normalizeOrderDetail(order: OrderDetailApi): Order {
  return {
    order_sn: order.order_sn,
    shop_id: order.shop_id,
    marketplace_status: order.marketplace_status as Order["marketplace_status"],
    shipping_status: order.shipping_status as Order["shipping_status"],
    wms_status: order.wms_status as WmsStatus,
    tracking_number: order.tracking_number || null,
    total_amount: order.total_amount,
    created_at: "",
    updated_at: "",
    raw_marketplace_payload: {},
    items: (order.items ?? []).map((item: OrderItemApi) => ({
      sku: item.sku,
      quantity: item.quantity,
      price: item.price
    }))
  };
}

export async function getOrders(params: {
  page: number;
  limit: number;
  search: string;
  marketplaceStatuses: Order["marketplace_status"][];
  shippingStatuses: Order["shipping_status"][];
  wmsStatuses: Order["wms_status"][];
  sortBy: string;
  sortDir: SortDirection;
}): Promise<OrdersResponse> {
  const {
    page,
    limit,
    search,
    marketplaceStatuses,
    shippingStatuses,
    wmsStatuses,
    sortBy,
    sortDir
  } = params;
  const query = new URLSearchParams({
    page: String(page),
    limit: String(limit),
    sort_by: sortBy,
    sort_dir: sortDir
  });

  if (search.trim()) {
    query.set("search", search.trim());
  }

  if (wmsStatuses.length > 0) {
    query.set("wms_status", wmsStatuses.join(","));
  }

  if (marketplaceStatuses.length > 0) {
    query.set("marketplace_status", marketplaceStatuses.join(","));
  }

  if (shippingStatuses.length > 0) {
    query.set("shipping_status", shippingStatuses.join(","));
  }

  const data = await apiFetch<OrderListDataApi>(
    `${ordersApiUrl}?${query.toString()}`
  );
  return {
    orders: (data.orders ?? []).map(normalizeOrderListItem),
    page: data.page,
    limit: data.limit,
    total: data.total,
    total_pages: data.total_pages
  };
}

export async function getOrder(orderSn: string): Promise<Order> {
  const data = await apiFetch<OrderDetailApi>(`${ordersApiUrl}/${orderSn}`);
  return normalizeOrderDetail(data);
}

export async function getLogisticChannels(): Promise<LogisticChannel[]> {
  return apiFetch<LogisticChannel[]>("/marketplace/logistic/channels");
}

export async function updateOrderStatus(
  orderSn: string,
  action: "pick" | "pack" | "ship",
  channelId?: string
): Promise<Order> {
  if (action === "ship") {
    if (!channelId) {
      throw new Error("Please select a logistic channel");
    }

    await apiFetch(`${ordersApiUrl}/${orderSn}/ship`, {
      method: "POST",
      body: JSON.stringify({ channel_id: channelId })
    });
  } else {
    await apiFetch(`${ordersApiUrl}/${orderSn}/${action}`, {
      method: "POST"
    });
  }

  return getOrder(orderSn);
}
