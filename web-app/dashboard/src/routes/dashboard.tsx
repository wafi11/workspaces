import { useLogout, useProfile } from '@/features/api';
import { Sidebar } from '@/features/layout/Sidebar'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard')({
  component: RouteComponent,
})

function RouteComponent() {
  const { mutate } = useLogout();
  const { data: profileData } = useProfile();
  return <div
        className="flex  w-full h-full"
        style={{ background: "var(--color-app-bg)" }}
      >
        <Sidebar
          role={profileData?.role}
          userEmail={profileData?.email}
          userName={profileData?.username}
          onLogout={mutate}
        />
       <Outlet />
      </div>
}
