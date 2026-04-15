import { useLogout, useProfile } from "@/features/api/auth";
import { Sidebar } from "@/features/layout/Sidebar";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/templates")({
  component: RouteComponent,
});

function RouteComponent() {
  const { mutate } = useLogout();
  const { data: profileData } = useProfile();

  return (
    <div
      className="flex  w-full h-full overflow-hidden"
      style={{ background: "var(--color-app-bg)" }}
    >
      <Sidebar
        role="admin"
        userEmail={profileData?.email}
        userName={profileData?.username}
        onLogout={mutate}
      />
      <Outlet />
    </div>
  );
}
