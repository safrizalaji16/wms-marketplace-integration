import { ArrowDownWideNarrow, Search } from "lucide-react";
import type { ChangeEvent } from "react";
import { useOrderStore } from "../../store/useOrderStore";
import { FilterPopover } from "../ui/FilterPopover";
import type { MarketplaceStatus, ShippingStatus, WmsStatus } from "../../types/order";

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

export function OrdersToolbar() {
  const {
    search,
    setSearch,
    marketplaceStatuses: selectedMarketplace,
    shippingStatuses: selectedShipping,
    wmsStatuses: selectedWms,
    sortDirection,
    setSortDirection,
    toggleMarketplaceStatus,
    toggleShippingStatus,
    toggleWmsStatus
  } = useOrderStore();

  function handleSearchChange(event: ChangeEvent<HTMLInputElement>) {
    setSearch(event.target.value);
  }

  return (
    <div className="space-y-4">
      <div className="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
        <label className="flex items-center gap-3 rounded-3xl border border-slate-200 bg-white px-4 py-3 shadow-card">
          <Search size={18} className="text-slate-400" />
          <input
            value={search}
            onChange={handleSearchChange}
            placeholder="Search here..."
            className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
          />
        </label>

        <div className="rounded-3xl border border-slate-200 bg-white px-4 py-3 shadow-card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.25em] text-slate-400">Sorting</p>
              <p className="mt-1 text-sm font-semibold text-slate-700">Update At</p>
            </div>
            <button
              type="button"
              onClick={() => setSortDirection(sortDirection === "desc" ? "asc" : "desc")}
              className="inline-flex items-center gap-2 rounded-2xl bg-primary-50 px-3 py-2 text-sm font-semibold text-primary-700"
            >
              <ArrowDownWideNarrow size={16} />
              {sortDirection === "desc" ? "Z - A" : "A - Z"}
            </button>
          </div>
        </div>
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <FilterPopover
          title="Marketplace"
          values={marketplaceStatuses}
          selected={selectedMarketplace}
          onToggle={toggleMarketplaceStatus}
        />
        <FilterPopover
          title="Shipping"
          values={shippingStatuses}
          selected={selectedShipping}
          onToggle={toggleShippingStatus}
        />
        <FilterPopover
          title="WMS"
          values={wmsStatuses}
          selected={selectedWms}
          onToggle={toggleWmsStatus}
        />
      </div>
    </div>
  );
}
