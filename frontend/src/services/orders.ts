import { apiFetch } from "./api";
import type {
  LogisticChannel,
  OrderDetailApi,
  OrderItemApi,
  OrderListDataApi,
  OrderListItemApi
} from "../types/api";
import type { Order, OrdersResponse, WmsStatus } from "../types/order";

let inMemoryOrders = [] as Order[];
const useMockApi = import.meta.env.VITE_USE_MOCK_API !== "false";
const ordersApiUrl = "/api/orders";

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function sortOrders(orders: Order[]) {
  return [...orders].sort(
    (left, right) =>
      new Date(right.updated_at).getTime() - new Date(left.updated_at).getTime()
  );
}

function nextWmsStatus(status: WmsStatus, action: "pick" | "pack" | "ship") {
  if (action === "pick" && status === "READY_TO_PICK") {
    return "PICKING";
  }

  if (action === "pack" && status === "PICKING") {
    return "PACKED";
  }

  if (action === "ship" && status === "PACKED") {
    return "SHIPPED";
  }

  throw new Error("Invalid action for current WMS status");
}

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
    items: order.items.map((item: OrderItemApi) => ({
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
}): Promise<OrdersResponse> {
  const {
    page,
    limit,
    search,
    marketplaceStatuses,
    shippingStatuses,
    wmsStatuses
  } = params;

  if (!useMockApi) {
    const query = new URLSearchParams({
      page: String(page),
      limit: String(limit)
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
      orders: data.orders.map(normalizeOrderListItem),
      page: data.page,
      limit: data.limit,
      total: data.total,
      total_pages: data.total_pages
    };
  }

  await sleep(350);
  const normalizedSearch = search.trim().toLowerCase();
  const filteredOrders = inMemoryOrders.filter((order: Order) => {
    const matchesSearch =
      normalizedSearch.length === 0 ||
      order.order_sn.toLowerCase().includes(normalizedSearch) ||
      (order.tracking_number ?? "").toLowerCase().includes(normalizedSearch);

    const matchesMarketplace =
      marketplaceStatuses.length === 0 ||
      marketplaceStatuses.includes(order.marketplace_status);

    const matchesShipping =
      shippingStatuses.length === 0 ||
      shippingStatuses.includes(order.shipping_status);

    const matchesWms =
      wmsStatuses.length === 0 || wmsStatuses.includes(order.wms_status);

    return matchesSearch && matchesMarketplace && matchesShipping && matchesWms;
  });
  const sortedOrders = sortOrders(filteredOrders);
  const start = (page - 1) * limit;
  const orders = sortedOrders.slice(start, start + limit);

  return {
    orders,
    page,
    limit,
    total: sortedOrders.length,
    total_pages: Math.ceil(sortedOrders.length / limit)
  };
}

export async function getOrder(orderSn: string): Promise<Order> {
  if (!useMockApi) {
    const data = await apiFetch<OrderDetailApi>(`${ordersApiUrl}/${orderSn}`);
    return normalizeOrderDetail(data);
  }

  await sleep(250);
  const order = inMemoryOrders.find((item: Order) => item.order_sn === orderSn);

  if (!order) {
    throw new Error("Order not found");
  }

  return order;
}

export async function getLogisticChannels(): Promise<LogisticChannel[]> {
  if (!useMockApi) {
    return apiFetch<LogisticChannel[]>("/marketplace/logistic/channels");
  }

  await sleep(200);
  return [
    { id: "reg", name: "Regular Delivery" },
    { id: "sameday", name: "Same Day" },
    { id: "instant", name: "Instant Courier" }
  ];
}

export async function updateOrderStatus(
  orderSn: string,
  action: "pick" | "pack" | "ship",
  channelId?: string
): Promise<Order> {
  if (!useMockApi) {
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

  await sleep(500);
  const targetIndex = inMemoryOrders.findIndex((item: Order) => item.order_sn === orderSn);

  if (targetIndex === -1) {
    throw new Error("Order not found");
  }

  const current = inMemoryOrders[targetIndex];
  const wms_status = nextWmsStatus(current.wms_status, action);

  const updatedOrder: Order = {
    ...current,
    wms_status,
    tracking_number:
      action === "ship"
        ? current.tracking_number ?? `TRK-${crypto.randomUUID().slice(0, 8)}`
        : current.tracking_number,
    shipping_status: action === "ship" ? "shipped" : current.shipping_status,
    updated_at: new Date().toISOString()
  };

  inMemoryOrders[targetIndex] = updatedOrder;
  return updatedOrder;
}
