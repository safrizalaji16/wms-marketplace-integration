import { ArrowDownRight, ArrowUpRight } from "lucide-react";
import { cn } from "../../lib/utils";

export function MetricCard({
  label,
  value,
  trend,
  trendDown
}: {
  label: string;
  value: string;
  trend: string;
  trendDown?: boolean;
}) {
  return (
    <div className="rounded-3xl border border-slate-200/80 bg-white px-5 py-4 shadow-card">
      <p className="text-sm font-medium text-slate-500">{label}</p>
      <p className="mt-4 text-3xl font-extrabold tracking-tight text-slate-900">{value}</p>
      <p
        className={cn(
          "mt-2 inline-flex items-center gap-1 text-sm font-semibold",
          trendDown ? "text-rose-500" : "text-emerald-500"
        )}
      >
        {trendDown ? <ArrowDownRight size={15} /> : <ArrowUpRight size={15} />}
        {trend}
      </p>
    </div>
  );
}
