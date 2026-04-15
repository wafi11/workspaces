import { WorkspaceDetails } from '@/features/components/workspaces/WorkspaceDetails'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/workspaces/$workspaceId')({
  component: RouteComponent,
})

function RouteComponent() {
  const { workspaceId } = Route.useParams()
  return <WorkspaceDetails id={workspaceId}/>
}
