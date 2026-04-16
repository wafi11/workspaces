import { useProfile } from "@/features/api";
import { MainContainer } from "@/features/layout/MainContainer";
import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import { EmptyState } from "@/features/pages/home/EmptyState";
import { SectionQuotaUser } from "@/features/pages/home/SectionCardQuotaUser";
import { SectionTerminal } from "@/features/pages/home/SectionTerminal";
import { SectionWorkspace } from "@/features/pages/home/SectionWorkspace";
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
        <SectionQuotaUser />
        <SectionWorkspace />
                <SectionTerminal terminal_url={`https://${profileData.terminal_url}`}/>
        
      </MainContainer>
   </>
  );
}
