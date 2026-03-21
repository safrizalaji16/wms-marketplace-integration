import { Search } from "lucide-react";
import type { ChangeEvent } from "react";
import { useOrderStore } from "../../store/useOrderStore";

export function OrdersToolbar() {
  const {
    search,
    setSearch,
  } = useOrderStore();

  function handleSearchChange(event: ChangeEvent<HTMLInputElement>) {
    setSearch(event.target.value);
  }

  return (
    <div className="lg:grid-cols-[1.2fr_0.8fr]">
      <label className="flex items-center rounded-3xl border border-slate-200 bg-white px-4 py-3 shadow-card">
        <Search size={18} className="text-slate-400" />
        <input
          value={search}
          onChange={handleSearchChange}
          placeholder="Search here..."
          className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
        />
      </label>
    </div>
  );
}
