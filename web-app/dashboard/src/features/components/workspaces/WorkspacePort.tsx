import { useCreateWorkspacePort, useGetWorkspacesPort } from "@/features/api/workspace-port"
import { useState } from "react"
import { ActionBtn } from "../ActionButton"

interface WorkspacePortProps {
    workspaceId : string
}export function WorkspacePort({ workspaceId }: WorkspacePortProps) {
  const { mutate: createPort } = useCreateWorkspacePort({ workspaceId })
  const { data: ports } = useGetWorkspacesPort({ workspaceId })
  const [selected, setSelected] = useState("")

  const allowedPorts = Array.from({ length: 11 }, (_, i) => 3000 + i)

  return (
    <section className="mx-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <p className="text-sm font-medium" style={{ color: "var(--color-sidebar-text-active)" }}>
            Published ports
          </p>
          <p className="text-xs" style={{ color: "var(--color-sidebar-text-muted)" }}>
            Expose ports 3000–3010 to the internet
          </p>
        </div>
      </div>

      <div className="flex gap-2 mb-4">
        <select
          value={selected}
          onChange={(e) => setSelected(e.target.value)}
          className="flex-1 text-sm rounded-md px-3 py-1.5"
          style={{ background: "var(--color-sidebar-surface)", border: "1px solid var(--color-sidebar-border)", color: "var(--color-sidebar-text)" }}
        >
          <option value="">Select port...</option>
          {allowedPorts.map((p) => (
            <option key={p}>{p}</option>
          ))}
        </select>
        <ActionBtn
          label="Publish"
          variant="default"
          onClick={() => { if (selected) createPort({ data  : {port: parseInt(selected)} }) }}
        />
      </div>

      <div className="flex flex-col gap-2">
        {ports?.length === 0 && (
          <p className="text-xs text-center py-6" style={{ color: "var(--color-sidebar-text-muted)" }}>
            No ports published yet
          </p>
        )}
        {ports?.map((p : any) => (
          <div key={p.port} className="flex items-center justify-between px-3 py-2 rounded-md"
            style={{ background: "var(--color-sidebar-bg)", border: "1px solid var(--color-sidebar-border)" }}>
            <div className="flex items-center gap-3">
              <span className="text-xs font-medium px-2 py-0.5 rounded-md"
                style={{ background: "var(--color-sidebar-surface)", color: "var(--color-sidebar-text-active)" }}>
                :{p.port}
              </span>
              <span className="text-xs" style={{ color: "var(--color-sidebar-text-muted)" }}>{p.sub_domain}</span>
            </div>
            <div className="flex gap-2">
              <ActionBtn label="Open" variant="default" onClick={() => window.open(`https://${p.sub_domain}`, "_blank")} />
              {/* <ActionBtn label="Remove" variant="danger" onClick={() => deletePort({ port: p.port })} /> */}
            </div>
          </div>
        ))}
      </div>
    </section>
  )
}