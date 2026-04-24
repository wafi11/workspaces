
import { ProfilePage } from '@/features/components/settings/profile/Profile'
import { MainContainer } from '@/features/layout/MainContainer'
import { SidebarSettings } from '@/features/layout/SidebarSettings'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/settings/profile')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <SidebarSettings />
      <MainContainer>
        <ProfilePage />
      </MainContainer>
    </>
  )
}
