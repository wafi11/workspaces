"use client";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { useGetWorkspace } from "@/features/services/workspaces/api";
import { useGetMetricsWorkspaces } from "@/features/services/workspaces/sse";
import { cn } from "@/lib/utils";
import { relativeTime } from "@/utils/relativeTime";
import { Circle, Clock, ExternalLink, User } from "lucide-react";
import { MetricCard } from "./WorkspaceMetrics";
import Link from "next/link";

interface WorkspaceDetailsProps {
  slug: string;
}

export function WorkspaceDetails({ slug }: WorkspaceDetailsProps) {
  const { data } = useGetWorkspace(slug);
  const workspace = data;
  const { byApp } = useGetMetricsWorkspaces();

  const statusColor =
    workspace?.status === "running"
      ? "text-green-400"
      : workspace?.status === "error"
        ? "text-red-400"
        : "text-yellow-400";

  const pods = [data?.name];

  return (
    <>
      <TopBarAdmin
        title={workspace?.name ?? "Workspace"}
        description={workspace?.id as string}
      />

      {/* Info bar */}
      <div className="flex items-center gap-6 px-5 py-3 border-b border-[#1a1a1a] text-[11px] font-mono">
        <div className="flex items-center gap-1.5">
          <Circle className={cn("w-2 h-2 fill-current", statusColor)} />
          <span className={cn("font-semibold", statusColor)}>
            {workspace?.status ?? "—"}
          </span>
        </div>
        <span className="text-[#2a2a2a]">|</span>
        <div className="flex items-center gap-1.5 text-[#555]">
          <User className="w-3 h-3" />
          <span className="text-[#777]">{workspace?.user_id ?? "—"}</span>
        </div>
        <span className="text-[#2a2a2a]">|</span>
        <div className="flex items-center gap-1.5 text-[#555]">
          <Clock className="w-3 h-3" />
          <span className="text-[#777]">
            {workspace?.updated_at
              ? relativeTime(workspace.updated_at).toLocaleString()
              : "—"}
          </span>
        </div>
        <span className="text-[#2a2a2a]">|</span>
        {workspace?.url && workspace.url !== "/" && (
          <Link
            href={`https://${workspace.url}`}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1 text-blue-500 hover:text-blue-400 transition-colors"
          >
            <ExternalLink className="w-3 h-3" />
            <span>open</span>
          </Link>
        )}
      </div>

      {/* Metrics */}
      <div className="px-5 py-4 border-b border-[#1a1a1a] w-full">
        <span className="text-[10px] font-mono text-[#444] uppercase tracking-widest mb-3 block">
          resource usage
        </span>
        <div className="flex w-full">
          {pods.map((app) =>
            byApp[app as string] ? (
              <MetricCard
                key={app}
                label={app as string}
                cpu={byApp[app as string].total_cpu_cores}
                memory={byApp[app as string].total_memory_mb}
              />
            ) : null
          )}
        </div>
      </div>
    </>
  );
}
