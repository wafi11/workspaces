import { cn } from "@/lib/utils"

interface EmptyStateProps {
  className?: string
  icon?: React.ReactNode
  title?: string
  description?: string
}

const DefaultIcon = () => (
    <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth={1.5}
        strokeLinecap="round"
        strokeLinejoin="round"
        style={{ color: "var(--color-sidebar-text)" }}
    >
      <path d="M3 7a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z" />
      <path d="M8 3v4M16 3v4" />
    </svg>
)

export function EmptyState({
     className,
     icon = <DefaultIcon />,
     title = "Workspace Not Existing",
     description = "Create Your Workspace And Let's Coding",
   }: EmptyStateProps) {
  return (
      <div
          className={cn("flex flex-col items-center justify-center gap-3 py-16 mx-4 rounded-lg", className)}
          style={{ border: "1px dashed var(--color-sidebar-border)" }}
      >
        <div
            className="w-10 h-10 rounded-lg flex items-center justify-center"
            style={{ background: "var(--color-sidebar-surface)" }}
        >
          {icon}
        </div>
        <div className="flex flex-col items-center gap-1">
        <span
            className="text-sm font-medium"
            style={{ color: "var(--color-sidebar-text-active)" }}
        >
          {title}
        </span>
          <span
              className="text-xs"
              style={{ color: "var(--color-sidebar-text-muted)" }}
          >
          {description}
        </span>
        </div>
      </div>
  )
}