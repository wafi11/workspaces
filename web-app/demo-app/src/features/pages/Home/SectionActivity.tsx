import { CardActivity } from "@/features/components/CardActivity";
import { activityStyle, recentActivity, recentWorkspaces, statusStyle } from "@/features/data";

export function SectionActivity() {
  return (
    <section className="grid grid-cols-2 gap-3 mx-4">
      {/* Recent Workspaces */}
      <CardActivity
        title="Recent Workspaces"
        action="View all →"
        body={
          <div>
            {recentWorkspaces.map((ws) => {
              const s = statusStyle[ws.status];
              return (
                <div
                  key={ws.name}
                  className="flex items-center gap-3 px-4 py-2.5 cursor-pointer transition-colors"
                  style={{ borderBottom: "0.5px solid var(--color-app-border)" }}
                  onMouseEnter={(e) => (e.currentTarget.style.background = "#111")}
                  onMouseLeave={(e) => (e.currentTarget.style.background = "transparent")}
                >
                  <div
                    className="w-1.5 h-1.5 rounded-full shrink-0"
                    style={{ background: s.dot }}
                  />
                  <div className="flex flex-col gap-0.5 flex-1 min-w-0">
                    <span className="text-xs truncate" style={{ color: "#ccc" }}>
                      {ws.name}
                    </span>
                    <span className="text-[11px] truncate" style={{ color: "var(--color-sidebar-text-muted)" }}>
                      {ws.template}
                    </span>
                  </div>
                  <span
                    className="text-[10px] px-2 py-0.5 rounded shrink-0"
                    style={{ color: s.text, cssText: s.badge } as React.CSSProperties}
                  >
                    {ws.status}
                  </span>
                  <span className="text-[11px] shrink-0" style={{ color: "var(--color-sidebar-text-muted)" }}>
                    {ws.uptime}
                  </span>
                </div>
              );
            })}
          </div>
        }
      />

      {/* Recent Activity */}
      <CardActivity
        title="Recent Activity"
        body={
          <div>
            {recentActivity.map((act, i) => {
              const s = activityStyle[act.type];
              return (
                <div
                  key={i}
                  className="flex items-center gap-3 px-4 py-2.5"
                  style={{ borderBottom: "0.5px solid var(--color-app-border)" }}
                >
                  <div
                    className="w-6 h-6 rounded-md flex items-center justify-center text-[11px] shrink-0"
                    style={{ background: s.bg, color: s.color }}
                  >
                    {s.icon}
                  </div>
                  <div className="flex flex-col gap-0.5 flex-1 min-w-0">
                    <span className="text-xs" style={{ color: "#888" }}>
                      <span style={{ color: "#ccc" }}>{act.text}</span> {act.sub}
                    </span>
                    <span className="text-[11px]" style={{ color: "var(--color-sidebar-text-muted)" }}>
                      {act.time}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        }
      />
    </section>
  );
}