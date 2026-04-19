import { FindAllPAT } from '@/features/components/settings/SectionFindAllPat'
import { MainContainer } from '@/features/layout/MainContainer'
import { SidebarSettings } from '@/features/layout/SidebarSettings'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/settings/pat')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
        <SidebarSettings />
        <MainContainer>
          <FindAllPAT />
        </MainContainer>
    </>
  )
}
