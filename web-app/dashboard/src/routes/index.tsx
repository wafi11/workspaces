import { useLogout, useProfile } from "@/features/api/auth";
import { SectionHeader } from "@/features/components/workspaces";
import { EmptyState } from "@/features/components/workspaces/EmptyState";
import { SectionQuotaUser } from "@/features/components/workspaces/SectionCardQuotaUser";
import { SectionTerminal } from "@/features/components/workspaces/SectionTerminal";
import { MainContainer } from "@/features/layout/MainContainer";
import { Sidebar } from "@/features/layout/Sidebar";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { mutate } = useLogout();
  const { data: profileData } = useProfile();

  if (!profileData){
    return <EmptyState />
  }

  return (
    <div
      className="flex  w-full h-full"
      style={{ background: "var(--color-app-bg)" }}
    >
      <Sidebar
        role={profileData?.role}
        userEmail={profileData?.email}
        userName={profileData?.username}
        onLogout={mutate}
      />
      <MainContainer>
        <SectionHeader />
        <SectionQuotaUser />
        {/* <SectionWorkspace /> */}
        <SectionTerminal terminal_url={`https://${profileData.terminal_url}`}/>
      </MainContainer>
    </div>
  );
}
