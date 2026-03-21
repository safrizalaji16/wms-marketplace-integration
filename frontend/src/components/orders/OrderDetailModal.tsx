import { useEffect, useState } from "react";
import type { ReactNode } from "react";
import { LoaderCircle, PackageCheck, PackageOpen, Truck, X } from "lucide-react";
import { formatCurrency, formatDate, sentenceCase } from "../../lib/utils";
import {
  useLogisticChannels,
  useOrderAction,
  useOrderDetail
} from "../../hooks/useOrders";
import { useAuthStore } from "../../store/useAuthStore";
import { useOrderStore } from "../../store/useOrderStore";
import type { OrderItem, WmsStatus } from "../../types/order";
import { StatusBadge } from "../ui/StatusBadge";
import type { OrderFilterState } from "../../store/useOrderStore";
import type { AuthState } from "../../store/useAuthStore";

const actionMap: Record<
  WmsStatus,
  { action?: "pick" | "pack" | "ship"; label: string; icon: typeof PackageOpen }
> = {
  READY_TO_PICK: { action: "pick", label: "Pickup", icon: PackageOpen },
  PICKING: { action: "pack", label: "Pack", icon: PackageCheck },
  PACKED: { action: "ship", label: "Ship", icon: Truck },
  SHIPPED: { label: "Completed", icon: Truck }
};

const actionRoles: Record<"pick" | "pack" | "ship", string[]> = {
  pick: ["picker", "superAdmin"],
  pack: ["packer", "superAdmin"],
  ship: ["warehouseAdmin", "superAdmin"]
};

export function OrderDetailModal() {
  const role = useAuthStore((state: AuthState) => state.role);
  const selectedOrderSn = useOrderStore(
    (state: OrderFilterState) => state.selectedOrderSn
  );
  const closeOrder = useOrderStore((state: OrderFilterState) => state.closeOrder);
  const {
    data: order,
    isLoading,
    isError: isOrderError,
    error: orderError
  } = useOrderDetail(selectedOrderSn);
  const [selectedChannelId, setSelectedChannelId] = useState("");
  const orderAction = useOrderAction();
  const {
    data: logisticChannels = [],
    isLoading: isChannelsLoading,
    isError: isChannelsError,
    error: logisticChannelsError
  } = useLogisticChannels(Boolean(selectedOrderSn && order?.wms_status === "PACKED"));

  useEffect(() => {
    setSelectedChannelId("");
    orderAction.reset();
  }, [selectedOrderSn]);

  function handleClose() {
    orderAction.reset();
    setSelectedChannelId("");
    closeOrder();
  }

  useEffect(() => {
    if (order?.wms_status !== "PACKED") {
      setSelectedChannelId("");
      return;
    }

    if (!selectedChannelId && logisticChannels[0]) {
      setSelectedChannelId(logisticChannels[0].id);
    }
  }, [logisticChannels, order?.wms_status, selectedChannelId]);

  if (!selectedOrderSn) {
    return null;
  }

  const config = order ? actionMap[order.wms_status] : null;
  const ActionIcon = config?.icon ?? PackageOpen;
  const isUnauthorizedForAction = Boolean(
    config?.action && !actionRoles[config.action].includes(role)
  );

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/35 p-4 backdrop-blur-sm">
      <div className="w-full max-w-2xl rounded-[32px] bg-white p-5 shadow-panel">
        <div className="flex items-start justify-between">
          <div>
            <p className="text-xs font-bold uppercase tracking-[0.25em] text-slate-400">Order Detail</p>
            <h3 className="mt-2 text-2xl font-extrabold tracking-tight text-slate-900">{selectedOrderSn}</h3>
          </div>
          <button
            type="button"
            onClick={handleClose}
            className="rounded-2xl bg-slate-100 p-2 text-slate-500 transition hover:bg-slate-200"
          >
            <X size={18} />
          </button>
        </div>

        {isLoading ? (
          <div className="flex min-h-64 items-center justify-center text-slate-500">
            <LoaderCircle className="animate-spin" size={22} />
          </div>
        ) : isOrderError || !order ? (
          <div className="flex min-h-64 items-center justify-center rounded-[28px] border border-rose-200 bg-rose-50 px-6 text-center text-rose-600">
            {orderError instanceof Error ? orderError.message : "Failed to load order detail"}
          </div>
        ) : (
          <div className="mt-5 space-y-4">
            <div className="grid gap-4 rounded-[28px] border border-slate-200 bg-slate-50/70 p-4 md:grid-cols-2">
              <DetailBlock label="Order SN" value={order.order_sn} />
              <DetailBlock label="Shop ID" value={order.shop_id} />
              <DetailBlock
                label="Marketplace Status"
                value={<StatusBadge value={order.marketplace_status} />}
              />
              <DetailBlock
                label="Shipping Status"
                value={<StatusBadge value={order.shipping_status} />}
              />
              <DetailBlock label="WMS Status" value={<StatusBadge value={order.wms_status} />} />
              <DetailBlock label="Tracking Number" value={order.tracking_number ?? "-"} />
              <DetailBlock label="Total Amount" value={formatCurrency(order.total_amount)} />
              <DetailBlock label="Created At" value={formatDate(order.created_at)} />
              <DetailBlock label="Updated At" value={formatDate(order.updated_at)} />
              <DetailBlock label="Allowed Action" value={sentenceCase(config?.label ?? "-")} />
            </div>

            <div className="overflow-hidden rounded-[28px] border border-slate-200">
              <table className="min-w-full text-left text-sm">
                <thead className="bg-slate-50 text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
                  <tr>
                    <th className="px-4 py-3">SKU</th>
                    <th className="px-4 py-3">Qty</th>
                    <th className="px-4 py-3">Price</th>
                  </tr>
                </thead>
                <tbody>
                  {order.items.map((item: OrderItem) => (
                    <tr key={item.sku} className="border-t border-slate-100">
                      <td className="px-4 py-3 font-semibold text-slate-700">{item.sku}</td>
                      <td className="px-4 py-3 text-slate-500">{item.quantity}</td>
                      <td className="px-4 py-3 text-slate-500">{formatCurrency(item.price)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {order.wms_status === "PACKED" ? (
              <div className="rounded-[28px] border border-slate-200 bg-slate-50/70 p-4">
                <div className="flex flex-col gap-2">
                  <label
                    htmlFor="logistic-channel"
                    className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400"
                  >
                    Logistic Channel
                  </label>
                  <select
                    id="logistic-channel"
                    value={selectedChannelId}
                    onChange={(event) => setSelectedChannelId(event.target.value)}
                    disabled={isChannelsLoading || logisticChannels.length === 0}
                    className="rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm font-medium text-slate-700 outline-none disabled:cursor-not-allowed disabled:bg-slate-100"
                  >
                    <option value="">
                      {isChannelsLoading
                        ? "Loading channels..."
                        : logisticChannels.length === 0
                          ? "No channel available"
                          : "Select logistic channel"}
                    </option>
                    {logisticChannels.map((channel) => (
                      <option key={channel.id} value={channel.id}>
                        {channel.name}
                      </option>
                    ))}
                  </select>
                  {isChannelsError ? (
                    <p className="text-sm text-rose-600">
                      {logisticChannelsError instanceof Error
                        ? logisticChannelsError.message
                        : "Failed to load logistic channels"}
                    </p>
                  ) : null}
                </div>
              </div>
            ) : null}

            <button
              type="button"
              disabled={
                !config?.action ||
                isUnauthorizedForAction ||
                orderAction.isPending ||
                (config.action === "ship" && !selectedChannelId)
              }
              onClick={() =>
                config?.action
                  ? orderAction.mutate({
                      orderSn: order.order_sn,
                      action: config.action,
                      channelId:
                        config.action === "ship" ? selectedChannelId : undefined
                    })
                  : undefined
              }
              className="flex w-full items-center justify-center gap-2 rounded-[22px] bg-primary-600 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-primary-500/25 transition hover:bg-primary-700 disabled:cursor-not-allowed disabled:bg-slate-300"
            >
              {orderAction.isPending ? (
                <LoaderCircle size={18} className="animate-spin" />
              ) : (
                <ActionIcon size={18} />
              )}
              {config?.label ?? "Completed"}
            </button>

            {orderAction.isError ? (
              <div className="rounded-[24px] border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
                {orderAction.error instanceof Error
                  ? orderAction.error.message
                  : "Failed to update order"}
              </div>
            ) : null}

            {isUnauthorizedForAction ? (
              <div className="rounded-[24px] border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-700">
                Your role is not allowed to perform this action.
              </div>
            ) : null}
          </div>
        )}
      </div>
    </div>
  );
}

function DetailBlock({
  label,
  value
}: {
  label: string;
  value: ReactNode;
}) {
  return (
    <div>
      <p className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">{label}</p>
      <div className="mt-2 text-sm font-semibold text-slate-700">{value}</div>
    </div>
  );
}
