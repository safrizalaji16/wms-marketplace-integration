import { BellDot, ChevronDown, LogOut } from "lucide-react";
import { LogoMark } from "../icons/LogoMark";
import { cn } from "../../lib/utils";
import { useAuthStore } from "../../store/useAuthStore";
import type { AuthState } from "../../store/useAuthStore";

const navItems = ["Inbound", "Outbound", "Inventory", "Settings"] as const;

export function Header() {
  const logout = useAuthStore((state: AuthState) => state.logout);
  const userName = useAuthStore((state: AuthState) => state.userName);

  return (
    <header className="panel overflow-hidden bg-gradient-to-r from-primary-600 via-primary-500 to-[#6D8BFF] text-white">
      <div className="flex flex-col gap-5 px-5 py-5 md:flex-row md:items-center md:justify-between md:px-8">
        <div className="flex items-center gap-4">
          <LogoMark />
          <div>
            <p className="text-2xl font-extrabold tracking-tight">WMSpaceIO</p>
          </div>
        </div>

        <nav className="flex flex-wrap items-center gap-2 rounded-full bg-white/10 p-1">
          {navItems.map((item: (typeof navItems)[number]) => (
            <button
              key={item}
              className={cn(
                "rounded-full px-4 py-2 text-sm font-semibold transition",
                item === "Outbound"
                  ? "bg-white text-primary-700 shadow-lg"
                  : "text-white/75 hover:bg-white/10 hover:text-white"
              )}
            >
              {item}
            </button>
          ))}
        </nav>

        <div className="flex items-center gap-3 self-end md:self-auto">
          <button className="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/15 backdrop-blur transition hover:bg-white/20">
            <BellDot size={19} />
          </button>
          <button className="flex items-center gap-3 rounded-2xl bg-white/15 px-2 py-1.5 backdrop-blur transition hover:bg-white/20">
            <div className="h-9 w-9 rounded-2xl bg-[linear-gradient(135deg,#FFD3A5,#FD6585)]" />
            <div className="text-left">
              <p className="text-sm font-bold capitalize">{userName || "Operator"}</p>
              <p className="text-xs text-white/70">Warehouse staff</p>
            </div>
            <ChevronDown size={18} />
          </button>
          <button
            type="button"
            onClick={logout}
            className="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/15 backdrop-blur transition hover:bg-white/20"
          >
            <LogOut size={18} />
          </button>
        </div>
      </div>
    </header>
  );
}
