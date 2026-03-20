import { Search } from "lucide-react";
import { cn, sentenceCase } from "../../lib/utils";

interface FilterPopoverProps<T extends string> {
  title: string;
  values: readonly T[];
  selected: T[];
  onToggle: (value: T) => void;
}

export function FilterPopover<T extends string>({
  title,
  values,
  selected,
  onToggle
}: FilterPopoverProps<T>) {
  return (
    <div className="rounded-3xl border border-slate-200 bg-white p-3 shadow-card">
      <div className="flex items-center gap-2 rounded-2xl px-3 py-2 text-slate-400">
        <span className="text-xs font-semibold uppercase tracking-[0.2em] text-slate-300">
          {title}
        </span>
      </div>
      <div className="mt-3 space-y-2">
        {values.map((value: T) => {
          const checked = selected.includes(value);

          return (
            <button
              key={value}
              type="button"
              onClick={() => onToggle(value)}
              className={cn(
                "flex w-full items-center gap-3 rounded-2xl border px-3 py-2 text-left text-sm font-medium transition",
                checked
                  ? "border-primary-300 bg-primary-50 text-primary-700"
                  : "border-transparent bg-slate-50 text-slate-500 hover:border-slate-200"
              )}
            >
              <span
                className={cn(
                  "flex h-4 w-4 items-center justify-center rounded border",
                  checked ? "border-primary-500 bg-primary-500 text-white" : "border-slate-300 bg-white"
                )}
              >
                {checked ? "✓" : ""}
              </span>
              {sentenceCase(value)}
            </button>
          );
        })}
      </div>
      <div className="mt-4 flex items-center justify-between px-1">
        <span className="rounded-full bg-primary-600 px-2.5 py-1 text-xs font-bold text-white">
          Save
        </span>
      </div>
    </div>
  );
}
