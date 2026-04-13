import { useGetWorkspacesUsers } from "@/features/api/workspace";
import { EmptyState } from "./EmptyState";
import { WorkspaceCard } from "./WorkspaceCard";
import { DialogCreateWorkspaces } from "@/features/components/workspaces/DialogCreateWorkspaces";

export function SectionWorkspace() {
  const { data: workspaceData, isLoading } = useGetWorkspacesUsers();

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <div
            key={i}
            className="h-44 rounded-lg animate-pulse"
            style={{ background: "var(--color-sidebar-surface)" }}
          />
        ))}
      </div>
    );
  }

  if (!workspaceData || workspaceData.length === 0) {
    return <EmptyState />;
  }

  return (
    <section className="space-y-4 mt-8">
      <div className="flex justify-between items-center gap-2">
        <div className="flex items-center gap-2">
          <h2
            className="text-xs font-semibold uppercase tracking-wider"
            style={{ color: "var(--color-sidebar-text)" }}
          >
            Your Workspaces
          </h2>
          <span
            className="text-[10px] font-mono px-1.5 py-0.5 rounded"
            style={{
              color: "var(--color-sidebar-text)",
              background: "var(--color-sidebar-surface)",
              border: "1px solid var(--color-sidebar-border)",
            }}
          >
            {workspaceData.length} instances
          </span>
        </div>
        <DialogCreateWorkspaces />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
        {workspaceData.map((ws) => (
          <WorkspaceCard ws={ws} key={ws.id} />
        ))}
      </div>
    </section>
  );
}
