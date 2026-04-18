import { useGetWorkspace } from "@/features/api"
import { MainContainer } from "@/features/layout/MainContainer"
import { TopbarAdmin } from "@/features/layout/TopbarDashboard"
import { useGetMetricsWorkspaces } from "@/hooks/useGetWorkspaceMatrics"
import { formatDate } from "@/utils/formatDate"
import { InfoRow } from "./InfoRow"
import { MetricCard } from "./MetricCard"
import { EmptyState } from "./EmptyState"

export type WorkspaceDetailsProps = {
    id : string
}

export function WorkspaceDetails({id} : WorkspaceDetailsProps){
  const {data}  = useGetWorkspace(id)
  const {byApp} = useGetMetricsWorkspaces()
  const pods = [data?.name];
  const envEntries = Object.entries(data?.env_vars ?? {});

  if (!data) {
    return <EmptyState />
  }

  return (
      <MainContainer>

        {/* Topbar */}
        <TopbarAdmin title={data.name}/>

        {/* Resource usage */}
        {pods && (
          <section className="px-5 py-4 border-b border-[#111]">
            <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase mb-3 block">
              resource usage
            </span>
            <div className="grid grid-cols-2 gap-2">
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
          </section>
        )}

        {/* Workspace info */}
        <section className="px-5 py-4 border-b border-[#111]">
          <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase mb-3 block">
            workspace info
          </span>
          <div className="bg-[#0d0d0d] border border-[#141414] rounded-sm overflow-hidden">
            <InfoRow label="template" value={data.template_name ?? "—"} />
            {/* <InfoRow label="url" value={`https://${data.url}`} isLink /> */}
            <InfoRow label="created" value={formatDate(data.created_at)} />
          </div>
        </section>
        {/* <WorkspacePort workspaceId={id}/> */}

        {/* Env vars */}
        <section className="px-5 py-4">
          <div className="flex items-center justify-between mb-3">
            <span className="text-[9px] tracking-[0.18em] text-sidebar-text-muted uppercase">env vars</span>
            <span className="text-[9px] text-sidebar-text-muted tracking-[0.08em]">
              {envEntries.length} variables
            </span>
          </div>
          <div className="bg-[#0d0d0d] border border-[#141414] rounded-sm overflow-hidden">
            <div className="grid px-3 py-1.5 border-b border-[#111] bg-[#0a0a0a]"
              style={{ gridTemplateColumns: "160px 1fr" }}>
              <span className="text-[9px] text-[#2d2d2d] uppercase tracking-[0.12em]">key</span>
              <span className="text-[9px] text-[#2d2d2d] uppercase tracking-[0.12em] text-right">value</span>
            </div>
            {envEntries.map(([key, val]) => (
              <div key={key}
                className="grid px-3 py-2 border-b border-[#0f0f0f] last:border-0 items-center"
                style={{ gridTemplateColumns: "160px 1fr" }}>
                <span className="text-[10px] text-[#71717a]">{key}</span>
                <span className="text-[10px] text-sidebar-text-muted text-right truncate">{val}</span>
              </div>
            ))}
            {envEntries.length === 0 && (
              <p className="px-3 py-2.5 text-[10px] text-sidebar-text-muted">no variables configured</p>
            )}
          </div>
        </section>
      </MainContainer>
    )
  }
