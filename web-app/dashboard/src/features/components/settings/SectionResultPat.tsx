import { Copy } from "lucide-react"

interface SectionCreatePatProps {
    newToken : string
    handleCopy : ()  => void
    copied : boolean
    setNewToken : (data : string | null)  => void
}

export function SectionResultPat({copied,handleCopy,newToken,setNewToken} : SectionCreatePatProps){
    return (
        <div
          className="mb-4 p-3 rounded"
          style={{
            background: "var(--color-sidebar-surface)",
            border: "1px solid var(--color-app-border)",
          }}
        >
          <p className="text-[10px] mb-2" style={{ color: "var(--color-sidebar-text)" }}>
            Copy your token now — it won't be shown again.
          </p>
          <div className="flex items-center gap-2">
            <code
              className="flex-1 text-[11px] truncate px-2 py-1.5 rounded"
              style={{
                background: "var(--color-app-bg)",
                border: "1px solid var(--color-app-border)",
                color: "var(--color-sidebar-text-active)",
              }}
            >
              {newToken}
            </code>
            <button
              onClick={handleCopy}
              className="p-1.5 rounded transition-colors"
              style={{
                background: "var(--color-app-border)",
                color: copied ? "#4ade80" : "var(--color-sidebar-text-active)",
              }}
            >
              <Copy size={12} />
            </button>
          </div>
          <button
            onClick={() => setNewToken(null)}
            className="text-[10px] mt-2 transition-colors"
            style={{ color: "var(--color-sidebar-text)" }}
          >
            dismiss
          </button>
        </div>
    )
}