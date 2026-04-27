import type { Workspace } from "@/types";
import { WorkspaceActions } from "./WorkspacesAction";

export function ListView({ workspaces }: { workspaces: Workspace[] }) {
  return (
    <div>
      {workspaces.map((ws) => (
        <div
          key={ws.name}
          className="flex items-center gap-3 px-4 py-2.5 cursor-pointer"
          style={{ borderBottom: "0.5px solid var(--color-app-border)" }}
          onMouseEnter={(e) => (e.currentTarget.style.background = "var(--color-app-surface)")}
          onMouseLeave={(e) => (e.currentTarget.style.background = "transparent")}
        >
          {/* <WorkspaceFilterStatus status={ws.status as any} /> */}   
        <div
        className="w-7 h-7 rounded-md flex items-center justify-center shrink-0 overflow-hidden"
        style={{
            background: "var(--color-accent-surface)",
            border: "0.5px solid var(--color-accent-border)",
        }}
        >
        <img
            src={ws.icon}
            alt={ws.name}
            className="w-4 h-4 object-contain"
        />
        </div>

          <div className="flex flex-col gap-0.5 flex-1 min-w-0">
            <span className="text-xs font-medium truncate" style={{ color: "var(--color-sidebar-text-active)" }}>
              {ws.name}
            </span>
            <span className="text-[11px] truncate" style={{ color: "var(--color-sidebar-text-muted)" }}>
              {ws.description}
            </span>
          </div>

          <span className="text-[11px] w-24 shrink-0" style={{ color: "var(--color-sidebar-text)" }}>
            {ws.template}
          </span>

          <span
            className="text-[11px] font-mono truncate w-48 shrink-0"
            style={{ color: "var(--color-sidebar-text-muted)" }}
          >
            {ws.url}
          </span>

          <span className="text-[11px] w-16 text-right shrink-0" style={{ color: "var(--color-sidebar-text-muted)" }}>
            {ws.uptime}
          </span>

          <WorkspaceActions status={ws.status} />
        </div>
      ))}
    </div>
  );
}

export function GridView({ workspaces }: { workspaces: Workspace[] }) {
  return (
    <div className="grid grid-cols-2 lg:grid-cols-3 gap-2.5 p-3">
      {workspaces.map((ws) => (
        <div
          key={ws.name}
          className="flex flex-col gap-3 rounded-lg p-3 cursor-pointer"
          style={{
            background: "var(--color-app-surface)",
            border: "0.5px solid var(--color-app-border)",
          }}
          onMouseEnter={(e) => (e.currentTarget.style.borderColor = "var(--color-accent-border)")}
          onMouseLeave={(e) => (e.currentTarget.style.borderColor = "var(--color-app-border)")}
        >
          <div className="flex items-center gap-2">
            <div
              className="w-7 h-7 rounded-md flex items-center justify-center shrink-0 overflow-hidden"
              style={{
                background: "var(--color-accent-surface)",
                border: "0.5px solid var(--color-accent-border)",
              }}
            >
               <img
            src={ws.icon}
            alt={ws.name}
            className="w-4 h-4 object-contain"
        />
            </div>
            <div className="flex flex-col gap-0.5 flex-1 min-w-0">
              <span className="text-xs font-medium truncate" style={{ color: "var(--color-sidebar-text-active)" }}>
                {ws.name}
              </span>
              <span className="text-[11px]" style={{ color: "var(--color-sidebar-text-muted)" }}>
                {ws.template}
              </span>
            </div>
            {/* <WorkspaceFilterStatus status={ws.status as any} showLabel /> */}
          </div>

          <span
            className="text-[11px] font-mono truncate"
            style={{ color: "var(--color-sidebar-text-muted)" }}
          >
            {ws.url}
          </span>

          <div className="flex items-center justify-between">
            <span className="text-[11px]" style={{ color: "var(--color-sidebar-text-muted)" }}>
              {ws.uptime}
            </span>
            <WorkspaceActions status={ws.status} />
          </div>
        </div>
      ))}
    </div>
  );
}
