import { useGetWorkspacesUsers } from "@/features/api/workspace";
import { EmptyState } from "./EmptyState";
import { WorkspaceCard } from "./WorkspaceCard";
import { ButtonCreate } from "@/features/components/ButtonCreate";

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
    return (<>
    <div className="flex w-full justify-between items-center gap-2 p-4">
        <div className="flex items-center gap-2">
          <h2
            className="text-xs font-semibold uppercase tracking-wider"
            style={{ color: "var(--color-sidebar-text)" }}
          >
            Your Workspaces
          </h2>
          
        </div>
        <ButtonCreate label="+ Create Workspace" to="workspaces"  />
      </div>

    <EmptyState className=""/>
    </>
  )
  }

  return (
    <section className="m-8 space-y-2">
      <div className="flex justify-between items-center gap-2">
          <h2
            className="text-xs font-semibold uppercase tracking-wider"
            style={{ color: "var(--color-sidebar-text)" }}
          >
            Your Workspaces
          </h2>          
        <ButtonCreate label="+ Create Workspace" to="workspaces"  />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
        {workspaceData.map((ws) => (
          <WorkspaceCard ws={ws} key={ws.id} />
        ))}
      </div>
    </section>
  );
}
