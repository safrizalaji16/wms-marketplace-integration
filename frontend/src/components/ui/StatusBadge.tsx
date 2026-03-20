import { cn, sentenceCase } from "../../lib/utils";

const colorMap: Record<string, string> = {
  delivered: "bg-emerald-100 text-emerald-700",
  shipping: "bg-amber-100 text-amber-700",
  processing: "bg-violet-100 text-violet-700",
  paid: "bg-sky-100 text-sky-700",
  cancelled: "bg-rose-100 text-rose-700",
  approved: "bg-teal-100 text-teal-700",
  shipped: "bg-orange-100 text-orange-700",
  awaiting_pickup: "bg-indigo-100 text-indigo-700",
  ready_to_pick: "bg-orange-100 text-orange-700",
  picking: "bg-cyan-100 text-cyan-700",
  packed: "bg-lime-100 text-lime-700"
};

export function StatusBadge({ value }: { value: string }) {
  return (
    <span className={cn("status-pill", colorMap[value.toLowerCase()] ?? "bg-slate-100 text-slate-700")}>
      {sentenceCase(value)}
    </span>
  );
}
