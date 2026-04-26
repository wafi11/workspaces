import { DataUserQuota } from "@/features/data/dataUserQuota";

const iconMap: Record<string, JSX.Element> = {
  Workspaces: (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="4.67" width="12" height="8.66" rx="1.33"/>
      <path d="M5.33 2v2.67M10.67 2v2.67"/>
    </svg>
  ),
  RAM: (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="5.33" width="12" height="6.67" rx="1"/>
      <path d="M5.33 5.33V4M8 5.33V4M10.67 5.33V4"/>
    </svg>
  ),
  CPU: (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <rect x="4" y="4" width="8" height="8" rx="1"/>
      <path d="M6 2v2M10 2v2M6 12v2M10 12v2M2 6h2M2 10h2M12 6h2M12 10h2"/>
    </svg>
  ),
  Storage: (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <ellipse cx="8" cy="5.33" rx="6" ry="2"/>
      <path d="M2 5.33v5.34C2 12.01 4.69 13.33 8 13.33s6-1.32 6-2.66V5.33"/>
    </svg>
  ),
};

function getBarColor(pct: number) {
  if (pct >= 90) return "var(--color-danger, #e24b4a)";
  if (pct >= 70) return "var(--color-warn, #f59e0b)";
  return "var(--color-accent)";
}

function getValueColor(pct: number) {
  if (pct >= 90) return "var(--color-danger, #e24b4a)";
  if (pct >= 70) return "var(--color-warn, #f59e0b)";
  return "var(--color-accent-text)";
}

export function StatCards() {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
      {DataUserQuota.map((item) => {
        const pct = Math.round((item.used / item.max) * 100);

        return (
          <div
            key={item.name}
            className="flex flex-col gap-3 rounded-lg p-3.5"
            style={{
              background: "var(--color-app-surface)",
              border: "0.5px solid var(--color-app-border)",
            }}
          >
            <div className="flex items-center justify-between">
              <div
                className="flex items-center gap-1.5 text-[11px]"
                style={{ color: "var(--color-sidebar-text)" }}
              >
                {iconMap[item.name]}
                {item.name}
              </div>
              <span
                className="text-[10px]"
                style={{ color: getValueColor(pct) }}
              >
                {pct}%
              </span>
            </div>

            <div className="flex flex-col gap-1">
              <div className="flex items-baseline gap-1">
                <span
                  className="text-lg font-medium leading-none"
                  style={{ color: getValueColor(pct) }}
                >
                  {item.used}
                  {item.unit && (
                    <span
                      className="text-xs ml-0.5"
                      style={{ color: "var(--color-sidebar-text-muted)" }}
                    >
                      {item.unit}
                    </span>
                  )}
                </span>
                <span
                  className="text-[11px]"
                  style={{ color: "var(--color-sidebar-text-muted)" }}
                >
                  / {item.max} {item.unit}
                </span>
              </div>

              <div
                className="h-1 rounded-full overflow-hidden"
                style={{ background: "var(--color-sidebar-border)" }}
              >
                <div
                  className="h-full rounded-full transition-all"
                  style={{ width: `${pct}%`, background: getBarColor(pct) }}
                />
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}