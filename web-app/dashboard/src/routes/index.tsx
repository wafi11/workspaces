import { useLogout, useProfile } from "@/features/api/auth";
import { Sidebar } from "@/features/layout/Sidebar";
import { SectionHeader } from "@/features/pages/home";
import { SectionQuotaUser } from "@/features/pages/home/SectionCardQuotaUser";
import { SectionWorkspace } from "@/features/pages/home/SectionWorkspace";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { mutate } = useLogout();
  const { data: profileData } = useProfile();

  return (
    <div
      className="flex  w-full h-full"
      style={{ background: "var(--color-app-bg)" }}
    >
      <Sidebar
        role="admin"
        userEmail={profileData?.email}
        userName={profileData?.username}
        onLogout={mutate}
      />
      <main className="flex flex-col space-y-4 p-4 flex-1 min-w-0 h-full">
        <SectionHeader />
        <SectionQuotaUser />
        <SectionWorkspace />
      </main>
    </div>
  );
}
