import { useGetWorkspaceCollaboration } from "@/features/api/workspace-collaboration"
import type { WorkspaceCollaborations } from "@/types"
import { Users } from "lucide-react"

export function SectionWorkspaceCollaborations() {
  const { data: workspaceCollaborations } = useGetWorkspaceCollaboration()

  if (!workspaceCollaborations?.length) return null

  return (
    <section className="px-4 mt-6">
      {/* Header */}
      <div className="flex items-center gap-2 mb-3">
        <Users size={13} style={{ color: "var(--color-sidebar-text)" }} />
        <span className="text-xs font-mono uppercase tracking-widest"
          style={{ color: "var(--color-sidebar-text)" }}>
          Shared with me
        </span>
        <span className="text-[10px] font-mono px-1.5 py-0.5 rounded"
          style={{
            background: "var(--color-sidebar-surface)",
            color: "var(--color-sidebar-text)",
            border: "1px solid var(--color-sidebar-border)"
          }}>
          {workspaceCollaborations.length}
        </span>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
        {workspaceCollaborations.map((collab) => (
          <CollabCard key={collab.workspace_id} collab={collab} />
        ))}
      </div>
    </section>
  )
}

function CollabCard({ collab }: { collab: WorkspaceCollaborations }) {
  return (
    <div
      className="flex flex-col gap-3 p-4 rounded-lg"
      style={{
        background: "var(--color-sidebar-bg)",
        border: "1px solid var(--color-sidebar-border)",
      }}
    >
      {/* Top */}
      <div className="flex items-center gap-3">
        <div className="w-8 h-8 rounded-md flex items-center justify-center shrink-0 text-base overflow-hidden"
          style={{ background: "var(--color-sidebar-surface)", border: "1px solid var(--color-sidebar-border)" }}>
          {collab.template_icon ? (
            <img src={collab.template_icon} alt={collab.workspace_name}
              className="w-full h-full object-cover rounded-md" />
          ) : "🖥"}
        </div>

        <div className="flex flex-col min-w-0 flex-1">
          <span className="text-sm font-medium truncate"
            style={{ color: "var(--color-sidebar-text-active)" }}>
            {collab.workspace_name}
          </span>
          <span className="text-[11px] truncate"
            style={{ color: "var(--color-sidebar-text)" }}>
            {collab.template_name}
          </span>
        </div>

        {/* Role badge */}
        <span className="text-[10px] font-mono px-2 py-0.5 rounded shrink-0"
          style={{
            background: "var(--color-sidebar-surface)",
            color: collab.role === "editor"
              ? "var(--color-primary)"
              : "var(--color-sidebar-text)",
            border: "1px solid var(--color-sidebar-border)"
          }}>
          {collab.role}
        </span>
      </div>

      {/* Footer */}
      <div className="flex items-center justify-between pt-3 text-[11px] font-mono"
        style={{ borderTop: "1px solid var(--color-sidebar-border)" }}>
        <span style={{ color: "var(--color-sidebar-text-muted)" }}>
          Joined {new Date(collab.invited_at).toLocaleDateString("en-GB", {
            day: "numeric", month: "short", year: "numeric"
          })}
        </span>
        <button
          className="px-2.5 py-1 rounded transition-colors"
          style={{
            background: "var(--color-sidebar-surface)",
            color: "var(--color-sidebar-text-active)",
            border: "1px solid var(--color-sidebar-border)"
          }}
          onMouseEnter={e => (e.currentTarget.style.borderColor = "var(--color-primary-hover)")}
          onMouseLeave={e => (e.currentTarget.style.borderColor = "var(--color-sidebar-border)")}
          onClick={() => window.open(`https://${collab.workspace_url}`, "_blank")}
        >
          Open
        </button>
      </div>
    </div>
  )
}