import { Trash2 } from "lucide-react"
import { formatDate, formatDay } from "@/utils/formatDate"
import type { PAT } from "@/types"
import { LoadingScreen } from "@/features/layout/loadingScreen"
import { EmptyState } from "../workspaces/EmptyState"

type Props = {
  pats?: PAT[]
  isLoading: boolean
  onDelete: (id: string) => void
}

export function PATList({ pats, isLoading, onDelete }: Props) {
  if (isLoading) return <LoadingScreen />
  if (!pats?.length) return <EmptyState title="No Tokens Yet" className="px-0 m-0" description="Add Your Personal Access Token"/>

  return (
    <div className="rounded border overflow-hidden">
      {pats.map((pat, i) => (
        <div
          key={pat.id}
          className="flex items-center justify-between px-4 py-3 border-b last:border-none"
        >
          <div>
            <p className="text-[12px]">{pat.name}</p>
            <p className="text-[10px]">
              {pat.last_used_at
                ? `last used ${formatDate(pat.last_used_at)}`
                : "never used"}
              {pat.expires_at && ` · expires ${formatDay(pat.expires_at)}`}
            </p>
          </div>

          <button
            onClick={() => onDelete(pat.id)}
            className="p-1.5 opacity-40 hover:opacity-100 text-red-500"
          >
            <Trash2 size={13} />
          </button>
        </div>
      ))}
    </div>
  )
}