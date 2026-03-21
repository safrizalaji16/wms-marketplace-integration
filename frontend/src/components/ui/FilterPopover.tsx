import { useEffect, useState } from "react";
import { ArrowDownUp, ListFilter } from "lucide-react";
import { cn, sentenceCase } from "../../lib/utils";

type SortDirection = "asc" | "desc";

interface FilterPopoverProps<T extends string> {
  title: string;
  values?: readonly T[];
  selected?: T[];
  onToggle?: (value: T) => void;
  sortDirection: SortDirection;
  onSortChange: (value: SortDirection) => void;
  onSave?: () => void;
  onReset?: () => void;
}

export function FilterPopover<T extends string>({
  title,
  values = [],
  selected = [],
  onToggle,
  sortDirection,
  onSortChange,
  onSave,
  onReset
}: FilterPopoverProps<T>) {
  const hasFilter = values.length > 0 && Boolean(onToggle);
  const [activePanel, setActivePanel] = useState<"sort" | "filter">(
    hasFilter ? "filter" : "sort"
  );

  useEffect(() => {
    setActivePanel(hasFilter ? "filter" : "sort");
  }, [hasFilter, title]);

  return (
    <div className="rounded-[24px] border border-slate-200 bg-white p-3 shadow-card">
      <div className="mb-3 px-1 text-xs font-semibold text-slate-400">{title}</div>
      <div className="flex gap-3">
        <div className="flex flex-col gap-2">
          <SideButton
            active={activePanel === "sort"}
            onClick={() => setActivePanel("sort")}
            icon={<ArrowDownUp size={16} />}
          />
          {hasFilter ? (
            <SideButton
              active={activePanel === "filter"}
              onClick={() => setActivePanel("filter")}
              icon={<ListFilter size={16} />}
            />
          ) : null}
        </div>

        <div className="min-w-0 flex-1">
          {activePanel === "sort" ? (
            <div className="space-y-2">
              <SortButton
                active={sortDirection === "asc"}
                label="A - Z"
                onClick={() => onSortChange("asc")}
              />
              <SortButton
                active={sortDirection === "desc"}
                label="Z - A"
                onClick={() => onSortChange("desc")}
              />
            </div>
          ) : (
            <div className="space-y-2">
              {values.map((value) => {
                const checked = selected.includes(value);

                return (
                  <button
                    key={value}
                    type="button"
                    onClick={() => onToggle?.(value)}
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
                        checked
                          ? "border-primary-500 bg-primary-500 text-white"
                          : "border-slate-300 bg-white"
                      )}
                    >
                      {checked ? "✓" : ""}
                    </span>
                    {sentenceCase(value)}
                  </button>
                );
              })}
            </div>
          )}
        </div>
      </div>
      <div className="mt-4 flex items-center justify-between px-1">
        <button
          type="button"
          onClick={onReset}
          className="text-xs font-semibold text-slate-400 transition hover:text-slate-600"
        >
          Reset
        </button>
        <button
          type="button"
          onClick={onSave}
          className="rounded-full bg-primary-600 px-2.5 py-1 text-xs font-bold text-white"
        >
          Save
        </button>
      </div>
    </div>
  );
}

function SideButton({
  active,
  onClick,
  icon
}: {
  active: boolean;
  onClick: () => void;
  icon: React.ReactNode;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "flex h-9 w-9 items-center justify-center rounded-xl border transition",
        active
          ? "border-primary-300 bg-primary-50 text-primary-600"
          : "border-slate-200 text-slate-400 hover:bg-slate-50"
      )}
    >
      {icon}
    </button>
  );
}

function SortButton({
  active,
  label,
  onClick
}: {
  active: boolean;
  label: string;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "flex w-full items-center rounded-2xl border px-3 py-2 text-left text-sm font-medium transition",
        active
          ? "border-primary-300 bg-primary-50 text-primary-700"
          : "border-slate-200 bg-white text-slate-500 hover:bg-slate-50"
      )}
    >
      {label}
    </button>
  );
}
