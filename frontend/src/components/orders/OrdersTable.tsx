import { useEffect, useState } from "react";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsUpDown,
  Funnel
} from "lucide-react";
import { cn, formatDate } from "../../lib/utils";
import { useOrderStore } from "../../store/useOrderStore";
import type {
  MarketplaceStatus,
  Order,
  ShippingStatus,
  WmsStatus
} from "../../types/order";
import { StatusBadge } from "../ui/StatusBadge";
import type { OrderFilterState } from "../../store/useOrderStore";
import { FilterPopover } from "../ui/FilterPopover";

const marketplaceStatuses = [
  "processing",
  "paid",
  "shipping",
  "delivered",
  "cancelled"
] as const satisfies readonly MarketplaceStatus[];

const shippingStatuses = [
  "label_created",
  "awaiting_pickup",
  "shipped",
  "delivered",
  "cancelled"
] as const satisfies readonly ShippingStatus[];

const wmsStatuses = [
  "READY_TO_PICK",
  "PICKING",
  "PACKED",
  "SHIPPED"
] as const satisfies readonly WmsStatus[];

type ActiveMenu =
  | "order-sn"
  | "marketplace"
  | "shipping"
  | "wms"
  | "tracking-number"
  | "update-at"
  | null;

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
  const [activeMenu, setActiveMenu] = useState<ActiveMenu>(null);
  const [draftMarketplace, setDraftMarketplace] = useState<MarketplaceStatus[]>([]);
  const [draftShipping, setDraftShipping] = useState<ShippingStatus[]>([]);
  const [draftWms, setDraftWms] = useState<WmsStatus[]>([]);
  const [draftSortDir, setDraftSortDir] = useState<"asc" | "desc">("desc");

  const openOrder = useOrderStore((state: OrderFilterState) => state.openOrder);
  const setPage = useOrderStore((state: OrderFilterState) => state.setPage);
  const setLimit = useOrderStore((state: OrderFilterState) => state.setLimit);
  const sortBy = useOrderStore((state: OrderFilterState) => state.sortBy);
  const sortDirection = useOrderStore((state: OrderFilterState) => state.sortDirection);
  const setSort = useOrderStore((state: OrderFilterState) => state.setSort);
  const selectedMarketplace = useOrderStore(
    (state: OrderFilterState) => state.marketplaceStatuses
  );
  const selectedShipping = useOrderStore(
    (state: OrderFilterState) => state.shippingStatuses
  );
  const selectedWms = useOrderStore((state: OrderFilterState) => state.wmsStatuses);
  const setMarketplaceStatuses = useOrderStore(
    (state: OrderFilterState) => state.setMarketplaceStatuses
  );
  const setShippingStatuses = useOrderStore(
    (state: OrderFilterState) => state.setShippingStatuses
  );
  const setWmsStatuses = useOrderStore(
    (state: OrderFilterState) => state.setWmsStatuses
  );

  useEffect(() => {
    if (activeMenu) {
      setDraftSortDir(sortDirection);
    }
  }, [activeMenu]);

  const startEntry = total === 0 ? 0 : (page - 1) * limit + 1;
  const endEntry = total === 0 ? 0 : Math.min(page * limit, total);
  const startPage = Math.max(1, Math.min(page - 2, Math.max(totalPages - 4, 1)));
  const visiblePages = Array.from(
    { length: Math.min(5, totalPages || 1) },
    (_, index: number) => startPage + index
  ).filter((pageNumber: number) => pageNumber >= 1 && pageNumber <= totalPages);


  function toggleDraftValue<T extends string>(values: T[], value: T) {
    return values.includes(value)
      ? values.filter((item: T) => item !== value)
      : [...values, value];
  }

  function openMenu(type: Exclude<ActiveMenu, null>) {
    setActiveMenu((current) => (current === type ? null : type));

    if (type === "marketplace") {
      setDraftMarketplace(selectedMarketplace);
    }

    if (type === "shipping") {
      setDraftShipping(selectedShipping);
    }

    if (type === "wms") {
      setDraftWms(selectedWms);
    }
  }

  const activeMenuPopover = (() => {
    if (!activeMenu) {
      return null;
    }

    if (activeMenu === "order-sn") {
      return (
        <PopoverWrap align="left" className="left-4 top-[72px]">
          <FilterPopover
            title="Order SN"
            sortDirection={draftSortDir}
            onSortChange={setDraftSortDir}
            onReset={() => {
              setSort("updated_at", "desc");
              setActiveMenu(null);
            }}
            onSave={() => {
              setSort("order_sn", draftSortDir);
              setActiveMenu(null);
            }}
          />
        </PopoverWrap>
      );
    }

    if (activeMenu === "marketplace") {
      return (
        <PopoverWrap align="left" className="left-[170px] top-[72px]">
          <FilterPopover
            title="Marketplace Status"
            values={marketplaceStatuses}
            selected={draftMarketplace}
            onToggle={(value) =>
              setDraftMarketplace((current) => toggleDraftValue(current, value))
            }
            sortDirection={draftSortDir}
            onSortChange={setDraftSortDir}
            onReset={() => {
              setDraftMarketplace([]);
              setSort("updated_at", "desc");
            }}
            onSave={() => {
              setMarketplaceStatuses(draftMarketplace);
              setSort("marketplace_status", draftSortDir);
              setActiveMenu(null);
            }}
          />
        </PopoverWrap>
      );
    }

    if (activeMenu === "shipping") {
      return (
        <PopoverWrap align="left" className="left-[385px] top-[72px]">
          <FilterPopover
            title="Shipping Status"
            values={shippingStatuses}
            selected={draftShipping}
            onToggle={(value) =>
              setDraftShipping((current) => toggleDraftValue(current, value))
            }
            sortDirection={draftSortDir}
            onSortChange={setDraftSortDir}
            onReset={() => {
              setDraftShipping([]);
              setSort("updated_at", "desc");
            }}
            onSave={() => {
              setShippingStatuses(draftShipping);
              setSort("shipping_status", draftSortDir);
              setActiveMenu(null);
            }}
          />
        </PopoverWrap>
      );
    }

    if (activeMenu === "wms") {
      return (
        <PopoverWrap align="left" className="left-[560px] top-[72px]">
          <FilterPopover
            title="WMS Status"
            values={wmsStatuses}
            selected={draftWms}
            onToggle={(value) =>
              setDraftWms((current) => toggleDraftValue(current, value))
            }
            sortDirection={draftSortDir}
            onSortChange={setDraftSortDir}
            onReset={() => {
              setDraftWms([]);
              setSort("updated_at", "desc");
            }}
            onSave={() => {
              setWmsStatuses(draftWms);
              setSort("wms_status", draftSortDir);
              setActiveMenu(null);
            }}
          />
        </PopoverWrap>
      );
    }

    if (activeMenu === "tracking-number") {
      return (
        <PopoverWrap align="right" className="right-[180px] top-[72px]">
          <FilterPopover
            title="Tracking Number"
            sortDirection={draftSortDir}
            onSortChange={setDraftSortDir}
            onReset={() => {
              setSort("updated_at", "desc");
              setActiveMenu(null);
            }}
            onSave={() => {
              setSort("tracking_number", draftSortDir);
              setActiveMenu(null);
            }}
          />
        </PopoverWrap>
      );
    }

    return (
      <PopoverWrap align="right" className="right-4 top-[72px]">
        <FilterPopover
          title="Update At"
          sortDirection={draftSortDir}
          onSortChange={setDraftSortDir}
          onReset={() => {
            setSort("updated_at", "desc");
            setActiveMenu(null);
          }}
          onSave={() => {
            setSort("updated_at", draftSortDir);
            setActiveMenu(null);
          }}
        />
      </PopoverWrap>
    );
  })();

  return (
    <div className="relative overflow-visible rounded-[32px] border border-slate-200 bg-white shadow-card">
      <div className="overflow-x-auto">
        <table className="min-w-full text-left">
          <thead className="border-b border-slate-200 bg-slate-50/80">
            <tr className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
              <th className="relative px-4 py-4">
                <HeaderButton
                  label="Order SN"
                  onClick={() => openMenu("order-sn")}
                  icon="sort"
                  active={sortBy === "order_sn"}
                />
              </th>
              <th className="relative px-4 py-4">
                <HeaderButton
                  label="Marketplace Status"
                  onClick={() => openMenu("marketplace")}
                  icon="filter"
                  active={selectedMarketplace.length > 0 || sortBy === "marketplace_status"}
                />
              </th>
              <th className="relative px-4 py-4">
                <HeaderButton
                  label="Shipping Status"
                  onClick={() => openMenu("shipping")}
                  icon="filter"
                  active={selectedShipping.length > 0 || sortBy === "shipping_status"}
                />
              </th>
              <th className="relative px-4 py-4">
                <HeaderButton
                  label="WMS Status"
                  onClick={() => openMenu("wms")}
                  icon="filter"
                  active={selectedWms.length > 0 || sortBy === "wms_status"}
                />
              </th>
              <th className="relative px-4 py-4">
                <HeaderButton
                  label="Tracking Number"
                  onClick={() => openMenu("tracking-number")}
                  icon="sort"
                  active={sortBy === "tracking_number"}
                />
              </th>
              <th className="relative px-4 py-4">
                <HeaderButton
                  label="Update At"
                  onClick={() => openMenu("update-at")}
                  icon="sort"
                  active={sortBy === "updated_at"}
                />
              </th>
              <th className="px-4 py-4">
                <div className="flex items-center gap-2">
                  <span>Action</span>
                </div>
              </th>
            </tr>
          </thead>
          <tbody>
            {orders.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-6 py-16">
                  <div className="flex min-h-56 flex-col items-center justify-center rounded-[28px] border border-dashed border-slate-200 bg-slate-50/70 text-center">
                    <p className="text-lg font-bold text-slate-700">No orders found</p>
                    <p className="mt-2 max-w-md text-sm leading-6 text-slate-500">
                      There are no outbound orders matching your current search or filters.
                      Try clearing some filters or changing the keyword.
                    </p>
                  </div>
                </td>
              </tr>
            ) : (
              orders.map((order: Order) => (
                <tr
                  key={order.order_sn}
                  className="border-b border-slate-100 text-sm text-slate-600 transition hover:bg-primary-50/40"
                >
                  <td className="px-4 py-4 font-semibold text-slate-700">{order.order_sn}</td>
                  <td className="px-4 py-4">
                    <StatusBadge value={order.marketplace_status} />
                  </td>
                  <td className="px-4 py-4">
                    <StatusBadge value={order.shipping_status} />
                  </td>
                  <td className="px-4 py-4">
                    <StatusBadge value={order.wms_status} />
                  </td>
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
              ))
            )}
          </tbody>
        </table>
      </div>

      {activeMenu && (
        <div className="fixed inset-0 z-[29]" onClick={() => setActiveMenu(null)} />
      )}
      {activeMenuPopover}

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
              className={`h-9 w-9 rounded-xl border text-sm font-semibold ${pageNumber === page
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

function PopoverWrap({
  children,
  align,
  className
}: {
  children: React.ReactNode;
  align: "left" | "right";
  className?: string;
}) {
  return (
    <div
      className={cn(
        "absolute z-30 w-80",
        align === "right" ? "right-0" : "left-0",
        className
      )}
    >
      {children}
    </div>
  );
}

function HeaderButton({
  label,
  onClick,
  icon,
  active
}: {
  label: string;
  onClick: () => void;
  icon: "sort" | "filter";
  active?: boolean;
}) {
  return (
    <div className="flex items-center gap-2">
      <span>{label}</span>
      <button
        type="button"
        onClick={onClick}
        className={cn(
          "rounded-lg p-1 transition",
          active
            ? "bg-primary-50 text-primary-600"
            : "text-slate-400 hover:bg-slate-100 hover:text-primary-600"
        )}
      >
        {icon === "filter" ? <Funnel size={14} /> : <ChevronsUpDown size={14} />}
      </button>
    </div>
  );
}
