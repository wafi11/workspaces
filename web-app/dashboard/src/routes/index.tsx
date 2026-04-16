import { useLogout, useProfile } from "@/features/api/auth";
import { MainContainer } from "@/features/layout/MainContainer";
import { Sidebar } from "@/features/layout/Sidebar";
import { SectionHeader } from "@/features/pages/home";
import { EmptyState } from "@/features/pages/home/EmptyState";
import { SectionQuotaUser } from "@/features/pages/home/SectionCardQuotaUser";
import { SectionTerminal } from "@/features/pages/home/SectionTerminal";
import { SectionWorkspace } from "@/features/pages/home/SectionWorkspace";
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
        <SectionWorkspace />
        <SectionTerminal terminal_url={profileData.terminal_url}/>
      </MainContainer>
    </div>
  );
}
