"use client";
import { PageContainer } from "@/components/layouts";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { useGetListWorkspace } from "@/features/services/workspaces/api";
import { WorkspaceCard } from "./WorkspaceCard";

export function WorkspacePage() {
  const { data } = useGetListWorkspace();
  const workspaceData = data?.data;

  return (
    <>
      <TopBarAdmin title="List Workspace" description="List Workspaces user">
        <></>
      </TopBarAdmin>
      {/* Section Content List Workspaces */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {workspaceData?.map((w) => (
          <WorkspaceCard key={w.id} workspace={w} />
        ))}
      </div>
    </>
  );
}
