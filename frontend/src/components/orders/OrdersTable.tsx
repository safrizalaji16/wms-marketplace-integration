import { ChevronLeft, ChevronRight, ChevronsUpDown } from "lucide-react";
import { formatDate } from "../../lib/utils";
import { useOrderStore } from "../../store/useOrderStore";
import type { Order } from "../../types/order";
import { StatusBadge } from "../ui/StatusBadge";
import type { OrderFilterState } from "../../store/useOrderStore";

export function OrdersTable({
  orders,
  page,
  limit,
  total,
  totalPages
}: {
  orders: Order[];
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}) {
  const openOrder = useOrderStore(
    (state: OrderFilterState) => state.openOrder
  );
  const setPage = useOrderStore((state: OrderFilterState) => state.setPage);
  const setLimit = useOrderStore((state: OrderFilterState) => state.setLimit);
  const startEntry = total === 0 ? 0 : (page - 1) * limit + 1;
  const endEntry = total === 0 ? 0 : Math.min(page * limit, total);
  const startPage = Math.max(1, Math.min(page - 2, Math.max(totalPages - 4, 1)));
  const visiblePages = Array.from(
    { length: Math.min(5, totalPages || 1) },
    (_, index: number) => startPage + index
  ).filter((pageNumber: number) => pageNumber >= 1 && pageNumber <= totalPages);

  return (
    <div className="overflow-hidden rounded-[32px] border border-slate-200 bg-white shadow-card">
      <div className="overflow-x-auto">
        <table className="min-w-full text-left">
          <thead className="border-b border-slate-200 bg-slate-50/80">
            <tr className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
              {["Order SN", "Marketplace Status", "Shipping Status", "WMS Status", "Tracking Number", "Update At", "Action"].map((label: string) => (
                <th key={label} className="px-4 py-4">
                  <div className="flex items-center gap-2">
                    <span>{label}</span>
                    {label !== "Action" ? <ChevronsUpDown size={14} /> : null}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {orders.map((order: Order) => (
              <tr key={order.order_sn} className="border-b border-slate-100 text-sm text-slate-600 transition hover:bg-primary-50/40">
                <td className="px-4 py-4 font-semibold text-slate-700">{order.order_sn}</td>
                <td className="px-4 py-4"><StatusBadge value={order.marketplace_status} /></td>
                <td className="px-4 py-4"><StatusBadge value={order.shipping_status} /></td>
                <td className="px-4 py-4"><StatusBadge value={order.wms_status} /></td>
                <td className="px-4 py-4">{order.tracking_number ?? "-"}</td>
                <td className="px-4 py-4">{formatDate(order.updated_at)}</td>
                <td className="px-4 py-4">
                  <button
                    type="button"
                    onClick={() => openOrder(order.order_sn)}
                    className="rounded-2xl bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-primary-500/20 transition hover:bg-primary-700"
                  >
                    Detail
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="flex flex-col gap-3 px-4 py-4 text-sm text-slate-500 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
          <p>
            Show <span className="font-semibold text-slate-700">{startEntry}</span> to{" "}
            <span className="font-semibold text-slate-700">{endEntry}</span> of{" "}
            <span className="font-semibold text-slate-700">{total}</span> entries
          </p>
          <label className="flex items-center gap-2">
            <span>Rows</span>
            <select
              value={limit}
              onChange={(event) => setLimit(Number(event.target.value))}
              className="rounded-xl border border-slate-200 bg-white px-2 py-1 text-sm text-slate-700 outline-none"
            >
              {[10, 20, 50].map((size: number) => (
                <option key={size} value={size}>
                  {size}
                </option>
              ))}
            </select>
          </label>
        </div>
        <div className="flex items-center gap-2">
          <button
            type="button"
            disabled={page <= 1}
            onClick={() => setPage(page - 1)}
            className="rounded-xl border border-slate-200 p-2 text-slate-500 disabled:cursor-not-allowed disabled:opacity-40"
          >
            <ChevronLeft size={16} />
          </button>
          {visiblePages.map((pageNumber: number) => (
            <button
              key={pageNumber}
              type="button"
              onClick={() => setPage(pageNumber)}
              className={`h-9 w-9 rounded-xl border text-sm font-semibold ${
                pageNumber === page
                  ? "border-primary-300 bg-primary-50 text-primary-700"
                  : "border-slate-200 text-slate-500"
              }`}
            >
              {pageNumber}
            </button>
          ))}
          <button
            type="button"
            disabled={page >= totalPages}
            onClick={() => setPage(page + 1)}
            className="rounded-xl border border-slate-200 p-2 text-slate-500 disabled:cursor-not-allowed disabled:opacity-40"
          >
            <ChevronRight size={16} />
          </button>
        </div>
      </div>
    </div>
  );
}
