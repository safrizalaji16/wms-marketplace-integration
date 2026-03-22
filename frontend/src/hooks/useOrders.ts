import type { Order, SortDirection } from "../types/order";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getLogisticChannels,
  getOrder,
  getOrders,
  updateOrderStatus
} from "../services/orders";

export function useOrders({
  enabled = true,
  page,
  limit,
  search,
  marketplaceStatuses,
  shippingStatuses,
  wmsStatuses,
  sortBy,
  sortDirection
}: {
  enabled?: boolean;
  page: number;
  limit: number;
  search: string;
  marketplaceStatuses: Order["marketplace_status"][];
  shippingStatuses: Order["shipping_status"][];
  wmsStatuses: Order["wms_status"][];
  sortBy: string;
  sortDirection: SortDirection;
}) {
  return useQuery({
    queryKey: [
      "orders",
      page,
      limit,
      search,
      marketplaceStatuses,
      shippingStatuses,
      wmsStatuses,
      sortBy,
      sortDirection
    ],
    queryFn: () =>
      getOrders({
        page,
        limit,
        search,
        marketplaceStatuses,
        shippingStatuses,
        wmsStatuses,
        sortBy,
        sortDir: sortDirection
      }),
    enabled
  });
}

export function useOrderDetail(orderSn: string | null) {
  return useQuery({
    queryKey: ["order", orderSn],
    queryFn: () => getOrder(orderSn!),
    enabled: Boolean(orderSn)
  });
}

export function useLogisticChannels(enabled = true) {
  return useQuery({
    queryKey: ["logistic-channels"],
    queryFn: getLogisticChannels,
    enabled
  });
}

export function useOrderAction() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      orderSn,
      action,
      channelId
    }: {
      orderSn: string;
      action: "pick" | "pack" | "ship";
      channelId?: string;
    }) => updateOrderStatus(orderSn, action, channelId),
    onSuccess: (order: Order) => {
      queryClient.setQueryData(["order", order.order_sn], order);
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    }
  });
}
