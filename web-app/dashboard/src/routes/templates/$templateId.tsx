import { TemplateDetails } from '@/features/components/templates/TemplateDetails'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/templates/$templateId')({
  component: RouteComponent,
})

function RouteComponent() {
  const { templateId } = Route.useParams()
  return <TemplateDetails templateId={templateId} />
}