import { create } from "zustand";
import type {
  MarketplaceStatus,
  ShippingStatus,
  SortDirection,
  WmsStatus
} from "../types/order";

export interface OrderFilterState {
  search: string;
  marketplaceStatuses: MarketplaceStatus[];
  shippingStatuses: ShippingStatus[];
  wmsStatuses: WmsStatus[];
  sortDirection: SortDirection;
  page: number;
  limit: number;
  selectedOrderSn: string | null;
  setSearch: (value: string) => void;
  setMarketplaceStatuses: (values: MarketplaceStatus[]) => void;
  setShippingStatuses: (values: ShippingStatus[]) => void;
  setWmsStatuses: (values: WmsStatus[]) => void;
  toggleMarketplaceStatus: (value: MarketplaceStatus) => void;
  toggleShippingStatus: (value: ShippingStatus) => void;
  toggleWmsStatus: (value: WmsStatus) => void;
  setSortDirection: (value: SortDirection) => void;
  setPage: (value: number) => void;
  setLimit: (value: number) => void;
  resetFilters: () => void;
  openOrder: (orderSn: string) => void;
  closeOrder: () => void;
}

const initialState = {
  search: "",
  marketplaceStatuses: [] as MarketplaceStatus[],
  shippingStatuses: [] as ShippingStatus[],
  wmsStatuses: [] as WmsStatus[],
  sortDirection: "desc" as SortDirection,
  page: 1,
  limit: 10,
  selectedOrderSn: null
};

function toggleValue<T>(values: T[], value: T) {
  return values.includes(value)
    ? values.filter((item: T) => item !== value)
    : [...values, value];
}

export const useOrderStore = create<OrderFilterState>((set) => ({
  ...initialState,
  setSearch: (search: string) => set({ search, page: 1 }),
  setMarketplaceStatuses: (marketplaceStatuses: MarketplaceStatus[]) =>
    set({ marketplaceStatuses, page: 1 }),
  setShippingStatuses: (shippingStatuses: ShippingStatus[]) =>
    set({ shippingStatuses, page: 1 }),
  setWmsStatuses: (wmsStatuses: WmsStatus[]) => set({ wmsStatuses, page: 1 }),
  toggleMarketplaceStatus: (value: MarketplaceStatus) =>
    set((state: OrderFilterState) => ({
      marketplaceStatuses: toggleValue(state.marketplaceStatuses, value),
      page: 1
    })),
  toggleShippingStatus: (value: ShippingStatus) =>
    set((state: OrderFilterState) => ({
      shippingStatuses: toggleValue(state.shippingStatuses, value),
      page: 1
    })),
  toggleWmsStatus: (value: WmsStatus) =>
    set((state: OrderFilterState) => ({
      wmsStatuses: toggleValue(state.wmsStatuses, value),
      page: 1
    })),
  setSortDirection: (sortDirection: SortDirection) => set({ sortDirection, page: 1 }),
  setPage: (page: number) => set({ page }),
  setLimit: (limit: number) => set({ limit, page: 1 }),
  resetFilters: () => set(initialState),
  openOrder: (selectedOrderSn: string) => set({ selectedOrderSn }),
  closeOrder: () => set({ selectedOrderSn: null })
}));
