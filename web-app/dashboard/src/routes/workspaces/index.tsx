import { useProfile } from "@/features/api";
import { EmptyState } from "@/features/components/workspaces/EmptyState";
import { SectionWorkspace } from "@/features/components/workspaces/SectionWorkspace";
import { SectionWorkspaceCollaborations } from "@/features/components/workspaces/SectionWorkspaceCollaborations";
import { MainContainer } from "@/features/layout/MainContainer";
import { TopbarAdmin } from "@/features/layout/TopbarDashboard";

import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/workspaces/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: profileData } = useProfile();
  if (!profileData){
      return <EmptyState />
    }

  
  return (
   <>
        <MainContainer >
          <TopbarAdmin title="Workspaces" />
          <SectionWorkspace />
          <SectionWorkspaceCollaborations />
      </MainContainer>
   </>
  );
}
