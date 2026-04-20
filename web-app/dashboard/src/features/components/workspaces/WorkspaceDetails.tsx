import { useGetWorkspace } from "@/features/api";
import { MainContainer } from "@/features/layout/MainContainer";
import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import { useGetMetricsWorkspaces } from "@/hooks/useGetWorkspaceMatrics";
import { formatDate } from "@/utils/formatDate";
import { ExternalLink, PanelRight, X } from "lucide-react";
import { useState } from "react";
import { EmptyState } from "./EmptyState";
import { InfoRow } from "./InfoRow";
import { MetricCard } from "./MetricCard";
import { WorkspacePort } from "./WorkspacePort";

export type WorkspaceDetailsProps = {
  id: string;
};

export function WorkspaceDetails({ id }: WorkspaceDetailsProps) {
  const { data } = useGetWorkspace(id);
  const [showPanel, setShowPanel] = useState(false);

  const { byApp } = useGetMetricsWorkspaces();
  const pods = [data?.name];
  const envEntries = Object.entries(data?.env_vars ?? {});
  console.log(byApp);
  if (!data) return <EmptyState />;

  return (
    <MainContainer>
      <TopbarAdmin title={data.name}>
        <a
          href={`https://${data.url}`}
          target="_blank"
          className="p-1.5 rounded hover:bg-[#1a1a1a] text-sidebar-text-active hover:text-white transition-colors"
        >
          <ExternalLink size={15} />
        </a>
        <button
          onClick={() => setShowPanel((v) => !v)}
          className="p-1.5 rounded hover:bg-[#1a1a1a] text-sidebar-text-active hover:text-white transition-colors"
        >
          <PanelRight size={15} />
        </button>
      </TopbarAdmin>

      <div className="flex w-full" style={{ height: "calc(100vh - 48px)" }}>
        {/* iframe */}
        <section className="flex-1 min-w-0">
          <iframe
            src={`https://${data.url}`}
            className="w-full h-full"
            allow="clipboard-read; clipboard-write"
          />
        </section>

        {/* Side panel */}
        {showPanel && (
          <aside className="w-64 shrink-0 border-l border-[#111] bg-[#080808] overflow-y-auto flex flex-col">
            {/* Panel header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-[#111]">
              <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase">
                workspace
              </span>
              <button
                onClick={() => setShowPanel(false)}
                className="text-sidebar-text-muted hover:text-white transition-colors"
              >
                <X size={13} />
              </button>
            </div>

            {/* Metrics */}
            <div className="px-4 py-3 border-b border-[#111]">
              <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase mb-2 block">
                resource usage
              </span>
              <div className="flex flex-col gap-2">
                {pods.map((app) =>
                  byApp[app as string] ? (
                    <MetricCard
                      key={app}
                      label={app as string}
                      cpu={byApp[app as string].total_cpu_cores}
                      memory={byApp[app as string].total_memory_mb}
                    />
                  ) : (
                    <p
                      key={app}
                      className="text-[10px] text-sidebar-text-muted"
                    >
                      no metrics yet
                    </p>
                  ),
                )}
              </div>
            </div>

            {/* Info */}
            <div className="px-4 py-3 border-b border-[#111]">
              <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase mb-2 block">
                info
              </span>
              <div className="bg-[#0d0d0d] border border-[#141414] rounded-sm overflow-hidden">
                <InfoRow label="template" value={data.template_name ?? "—"} />
                <InfoRow label="created" value={formatDate(data.created_at)} />
              </div>
            </div>
            <WorkspacePort workspaceId={id} />

            {/* Env vars */}
            <div className="px-4 py-3">
              <div className="flex items-center justify-between mb-2">
                <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase">
                  env vars
                </span>
                <span className="text-[9px] text-sidebar-text-muted">
                  {envEntries.length}
                </span>
              </div>
              <div className="bg-[#0d0d0d] border border-[#141414] rounded-sm overflow-hidden">
                {envEntries.length === 0 ? (
                  <p className="px-3 py-2.5 text-[10px] text-sidebar-text-muted">
                    no variables configured
                  </p>
                ) : (
                  envEntries.map(([key, val]) => (
                    <div
                      key={key}
                      className="px-3 py-2 border-b border-[#0f0f0f] last:border-0"
                    >
                      <p className="text-[10px] text-[#71717a]">{key}</p>
                      <p className="text-[10px] text-sidebar-text-muted truncate">
                        {val}
                      </p>
                    </div>
                  ))
                )}
              </div>
            </div>
          </aside>
        )}
      </div>
    </MainContainer>
  );
}
