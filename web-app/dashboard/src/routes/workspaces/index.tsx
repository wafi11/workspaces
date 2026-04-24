import { SectionWorkspace } from "@/features/components/workspaces/SectionWorkspace";
import { SectionWorkspaceCollaborations } from "@/features/components/workspaces/SectionWorkspaceCollaborations";
import { MainContainer } from "@/features/layout/MainContainer";
import { TopbarAdmin } from "@/features/layout/TopbarDashboard";

import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/workspaces/")({
  component: RouteComponent,
});

function RouteComponent() {
  
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
