import { ADMIN_ROLE } from '@/constants'
import { useProfile } from '@/features/api'
import { CreateTemplates } from '@/features/components/templates/CreateTemplates'
import { LoadingScreen } from '@/features/layout/loadingScreen'
import { createFileRoute, useRouter } from '@tanstack/react-router'

export const Route = createFileRoute('/templates/create')({
  component: RouteComponent,
})

function RouteComponent() {
  const {navigate}  = useRouter()
  const {data : profileData} = useProfile()
  if (!profileData){
    return <LoadingScreen />
  }

  if (profileData.role !== ADMIN_ROLE){
    return navigate({to : "/"})
  }

  return <CreateTemplates />
}
