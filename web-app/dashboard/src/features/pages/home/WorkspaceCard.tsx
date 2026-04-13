import { useUpdateStatusWorkspace } from "@/features/api/workspace";
import { ActionBtn } from "@/features/components/ActionButton";
import { StatusBadge } from "@/features/components/statusBadge";
import type { WorkspaceSessions } from "@/types";
import { formatDate } from "@/utils/formatDate";

export function WorkspaceCard({ ws }: { ws: WorkspaceSessions }) {
  const isRunning = ws.status === "running";
  const isPaused = ws.status === "paused";
  const isStopped = ws.status === "stopped";
  const isPending = ws.status === "pending";

  const { mutate: startWs } = useUpdateStatusWorkspace(ws.id, "start");
  const { mutate: stopWs } = useUpdateStatusWorkspace(ws.id, "stop");
  const { mutate: resumeWs } = useUpdateStatusWorkspace(ws.id, "resumed");
  const { mutate: pauseWs } = useUpdateStatusWorkspace(ws.id, "paused");

  return (
    <div
      className="flex flex-col gap-4 p-4 rounded-lg transition-colors"
      style={{
        background: "var(--color-sidebar-bg)",
        border: "1px solid var(--color-sidebar-border)",
      }}
    >
      {/* Header: icon + name + status */}
      <div className="flex items-center gap-3">
        <div className="w-9 h-9 rounded-md flex items-center justify-center shrink-0 select-none overflow-hidden">
          {ws.icon ? (
            <>
              <img
                src={ws.icon}
                alt={ws.name}
                className="w-full h-full object-cover rounded-md"
                onError={(e) => {
                  e.currentTarget.style.display = "none";
                  e.currentTarget.nextElementSibling?.removeAttribute("hidden");
                }}
              />
              <span hidden className="text-lg">
                🖥
              </span>
            </>
          ) : (
            <span className="text-lg">🖥</span>
          )}
        </div>
        <div className="flex flex-col min-w-0 flex-1">
          <span
            className="text-sm font-semibold truncate"
            style={{ color: "var(--color-sidebar-text-active)" }}
          >
            {ws.name}
          </span>
        </div>
        <StatusBadge status={ws.status} />
      </div>

      {/* Timestamps */}
      <div
        className="grid grid-cols-2 gap-2 text-[11px] pt-3"
        style={{ borderTop: "1px solid var(--color-sidebar-border)" }}
      >
        <div className="flex flex-col gap-0.5">
          <span style={{ color: "var(--color-sidebar-text-muted)" }}>
            Expires
          </span>
          <span
            style={{
              color: ws.expires_at
                ? "var(--color-sidebar-text)"
                : "var(--color-sidebar-text-muted)",
            }}
          >
            {formatDate(ws.expires_at)}
          </span>
        </div>
        <div className="flex flex-col gap-0.5">
          <span style={{ color: "var(--color-sidebar-text-muted)" }}>
            Next start
          </span>
          <span
            style={{
              color: ws.next_start_at
                ? "var(--color-sidebar-text)"
                : "var(--color-sidebar-text-muted)",
            }}
          >
            {formatDate(ws.next_start_at)}
          </span>
        </div>
      </div>

      {/* Actions */}
      <div className="flex gap-2">
        {isRunning && (
          <>
            <ActionBtn label="Pause" variant="warn" onClick={() => pauseWs()} />
            <ActionBtn label="Stop" variant="danger" onClick={() => stopWs()} />
            <ActionBtn
              label="Open"
              variant="default"
              onClick={() => window.open(ws.url, "_blank")}
            />
          </>
        )}
        {isPaused && (
          <>
            <ActionBtn
              label="Resume"
              variant="default"
              onClick={() => resumeWs()}
            />
            <ActionBtn label="Stop" variant="danger" onClick={() => stopWs()} />
          </>
        )}
        {isStopped && (
          <ActionBtn
            label="Start"
            variant="default"
            onClick={() => startWs()}
          />
        )}
        {isPending && (
          <ActionBtn
            label="Start"
            variant="default"
            onClick={() => startWs()}
          />
        )}
      </div>
    </div>
  );
}
