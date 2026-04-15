import { TemplateDetailsProvider } from "./TemplateDetailsContext"
import { TemplateDetailsContent } from "./TemplateDetilsContent"

interface TemplateDetailsProps {
  templateId: string
}

export function TemplateDetails({ templateId }: TemplateDetailsProps) {
  return (
    <TemplateDetailsProvider templateId={templateId}>
      <TemplateDetailsContent />
    </TemplateDetailsProvider>
  )
}