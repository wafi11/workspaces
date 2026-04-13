"use client";

import { Badge } from "@/components/ui/badge";
import { useCountdown } from "@/features/hooks/workspaces/useCountdown";
import { WorkspaceSessions } from "@/types/workspaces";
import { Layers } from "lucide-react";
import { useRouter } from "next/navigation";
import { DialogCreateWorkspaces } from "../dashboard/workspaces/DialogCreateWorkspaces";
import { statusConfig } from "../dashboard/workspaces/Status";
import { useUpdateStatusWorkspace } from "@/features/services/workspaces/api";
import { Button } from "@/components/ui/button";
import { useWorkspaceSocket } from "@/hooks";
import { useQueryClient } from "@tanstack/react-query";

interface SectionWorkspacesProps {
  data?: WorkspaceSessions[];
}

export function SectionWorkspaces({ data }: SectionWorkspacesProps) {
  if (!data || data.length === 0) {
    return (
      <div className="p-8 text-center border-2 border-dashed rounded-lg">
        <p className="text-muted-foreground text-sm">
          No active workspaces found.
        </p>
        <DialogCreateWorkspaces title="Create Your Workspaces" />
      </div>
    );
  }

  return (
    <section className="space-y-4 mt-8">
      <div className="flex justify-between items-center gap-2">
        <div className="flex items-center gap-2">
          <Layers className="w-5 h-5 text-primary" />
          <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
            Your Workspaces
          </h2>
          <Badge variant="secondary" className="ml-2 font-mono text-[10px]">
            {data.length} Instances
          </Badge>
        </div>
        <DialogCreateWorkspaces className="mt-0 p-2 text-sm" title="Create" />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {data.map((ws) => (
          <SessionCard ws={ws} key={ws.id} />
        ))}
      </div>
    </section>
  );
}
export function SessionCard({ ws }: { ws: WorkspaceSessions }) {
  const { push } = useRouter();
  const { mutate: startWs, isPending: isStarting } = useUpdateStatusWorkspace(
    ws.id,
    "start"
  );
  const { mutate: stopWs, isPending: isStopping } = useUpdateStatusWorkspace(
    ws.id,
    "stop"
  );

  const s = statusConfig[ws.status] ?? statusConfig.stopped;
  const timeLeft = useCountdown(ws.expires_at);
  const timeNext = useCountdown(ws.next_start_at);

  const isInCooldown =
    ws.next_start_at && new Date(ws.next_start_at) > new Date();

  return (
    <div
      key={ws.id}
      className="overflow-hidden p-2 bg-muted/30 border-border/50 hover:bg-muted/50 transition-colors"
    >
      <div className="flex flex-col gap-2 p-3">
        {/* Header */}
        <div
          className="flex items-center justify-between gap-2 cursor-pointer"
          onClick={() => push(`/profile/workspaces/${ws.id}`)}
        >
          <div className="flex items-center gap-2 min-w-0">
            <div className="p-1.5 shrink-0 bg-background rounded-md border border-border/50">
              <img src={ws.icon} className="size-4" alt={ws.name} />
            </div>
            <h3 className="text-sm font-medium leading-none truncate">
              {ws.name}
            </h3>
          </div>
          <div className="flex items-center gap-1.5 shrink-0">
            <div className="relative flex h-2 w-2">
              <span
                className={`animate-ping absolute inline-flex h-full w-full rounded-full ${s.ping} opacity-75`}
              />
              <span
                className={`relative inline-flex rounded-full h-2 w-2 ${s.dot}`}
              />
            </div>
            <span
              className={`text-xs font-medium uppercase tracking-tighter ${s.color}`}
            >
              {ws.status}
            </span>
          </div>
        </div>

        {/* Timer */}
        {ws.expires_at && ws.status === "running" && (
          <p className="text-[10px] text-muted-foreground font-mono border-t border-border/40 pt-1.5">
            ⏱ {timeLeft}
          </p>
        )}
        {ws.next_start_at && ws.status === "stopped" && isInCooldown && (
          <p className="text-[10px] text-muted-foreground font-mono border-t border-border/40 pt-1.5">
            🕐 available in {timeNext}
          </p>
        )}

        {/* Actions */}
        <div className="flex gap-2 border-t border-border/40 pt-2">
          {ws.status === "stopped" && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                startWs();
              }}
              disabled={isStarting || !!isInCooldown}
              className="flex-1 text-[11px] font-mono py-1 px-2 rounded bg-green-500/10 text-green-400 hover:bg-green-500/20 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              {isStarting ? "starting..." : "start"}
            </button>
          )}
          {ws.status === "running" && (
            <Button
              onClick={(e) => {
                e.stopPropagation();
                stopWs();
              }}
              disabled={isStopping}
              className="flex-1 text-[11px] font-mono py-1 px-2 rounded bg-red-500/10 text-red-400 hover:bg-red-500/20 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              {isStopping ? "stopping..." : "stop"}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
