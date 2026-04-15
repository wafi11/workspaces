import { useCreateWorkspaces, useProfile } from '@/features/api'
import { FormWorkspaces } from '@/features/components/workspaces/FormWorkspace'
import { DeployLog } from '@/features/components/workspaces/WorkspaceLogs'
import { TopbarAdmin } from '@/features/layout/TopbarDashboard'
import type { WorkspaceRequest } from '@/types'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { z } from 'zod'

const templateSchema = z.object({
    templateId: z.string().optional(),
})


export const Route = createFileRoute('/workspaces/create')({
    validateSearch: (search) => templateSchema.parse(search),
  component: RouteComponent,
})

function RouteComponent() {
  const { templateId } = Route.useSearch()
  const { mutate } = useCreateWorkspaces()
  const {data : profile} = useProfile()
  const [workspaceId, setWorkspaceId] = useState<string | undefined>()

  const handleSuccess = (ws: WorkspaceRequest) => {
    mutate(ws, {
      onSuccess: (res) => {
        console.log(res)
        if (res.workspace.id) {
          setWorkspaceId(profile?.id) }
        }
    })
  }

  return (
    <main className="flex flex-col w-full h-screen">
      <TopbarAdmin title="Create Workspace" />
      <div className="grid grid-cols-2 gap-4 p-4 flex-1 overflow-y-auto">
        <FormWorkspaces onSuccess={handleSuccess} templateId={templateId}/>
        <DeployLog workspaceId={workspaceId} />
      </div>
    </main>
  )
}
