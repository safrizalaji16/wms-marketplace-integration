import { LoginScreen } from "./components/auth/LoginScreen";
import { AlertCircle, LoaderCircle } from "lucide-react";
import { Header } from "./components/layout/Header";
import { OrderDetailModal } from "./components/orders/OrderDetailModal";
import { OrdersTable } from "./components/orders/OrdersTable";
import { OrdersToolbar } from "./components/orders/OrdersToolbar";
import { MetricCard } from "./components/ui/MetricCard";
import { useOrders } from "./hooks/useOrders";
import { useAuthStore } from "./store/useAuthStore";
import { useOrderStore } from "./store/useOrderStore";
import { sentenceCase } from "./lib/utils";
import type { Order } from "./types/order";

function App() {
  const isAuthenticated = useAuthStore((state: { isAuthenticated: boolean }) => state.isAuthenticated);
  const {
    search,
    marketplaceStatuses,
    shippingStatuses,
    wmsStatuses,
    sortDirection,
    page,
    limit
  } = useOrderStore();
  const { data, isLoading, isError, error } = useOrders({
    enabled: isAuthenticated,
    page,
    limit,
    search,
    marketplaceStatuses,
    shippingStatuses,
    wmsStatuses
  });

  const orders = data?.orders ?? [];
  const sortedOrders = orders
    .sort((left: Order, right: Order) => {
      const first = new Date(left.updated_at).getTime();
      const second = new Date(right.updated_at).getTime();
      return sortDirection === "desc" ? second - first : first - second;
    });

  const cancelledCount = (data?.orders ?? []).filter(
    (order: Order) => order.marketplace_status === "cancelled"
  ).length;

  if (!isAuthenticated) {
    return <LoginScreen />;
  }

  return (
    <main className="relative min-h-screen overflow-hidden bg-halo px-4 py-5 text-ink md:px-6">
      <div className="mx-auto flex max-w-7xl flex-col gap-6">
        <Header />

        <section className="panel px-5 py-6 md:px-8 md:py-8">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <h1 className="mt-3 text-4xl font-extrabold tracking-tight text-slate-900">
                Outbound
              </h1>
              <p className="mt-3 max-w-2xl text-sm leading-6 text-slate-500">
                Manage all outbound process
              </p>
            </div>
            <div className="rounded-[28px] border border-primary-100 bg-primary-50 px-4 py-3 text-sm text-primary-700">
              Focused on <span className="font-extrabold">{data?.total ?? sortedOrders.length}</span> visible orders
            </div>
          </div>

          <div className="mt-8 grid gap-4 md:grid-cols-3">
            <MetricCard label="Total Order" value={`${data?.total ?? sortedOrders.length}`} trend="12% this month" />
            <MetricCard label="Cancelled" value={`${cancelledCount}`} trend="5% this month" trendDown />
            <MetricCard
              label="Ready To Pick"
              value={`${(data?.orders ?? []).filter((order: Order) => order.wms_status === "READY_TO_PICK").length}`}
              trend={`${sentenceCase("ready_to_pick")} queue`}
            />
          </div>

          <div className="mt-8">
            <OrdersToolbar />
          </div>

          <div className="mt-8">
            {isLoading ? (
              <div className="flex min-h-64 items-center justify-center rounded-[32px] border border-slate-200 bg-white">
                <LoaderCircle className="animate-spin text-primary-600" size={26} />
              </div>
            ) : isError ? (
              <div className="flex min-h-64 flex-col items-center justify-center gap-3 rounded-[32px] border border-rose-200 bg-rose-50 text-rose-600">
                <AlertCircle size={22} />
                <p className="font-semibold">
                  {error instanceof Error ? error.message : "Failed to load orders."}
                </p>
              </div>
            ) : (
              <OrdersTable
                orders={sortedOrders}
                page={page}
                limit={limit}
                total={data?.total ?? sortedOrders.length}
                totalPages={data?.total_pages ?? 1}
              />
            )}
          </div>
        </section>
      </div>

      <OrderDetailModal />
    </main>
  );
}

export default App;
