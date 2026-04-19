import { useCreatePat, useDeletePat, useFindAllPat } from "@/features/api/pat"
import { useState } from "react"
import { Plus } from "lucide-react"
import { SectionResultPat } from "./SectionResultPat"
import { SectionFormPat } from "./SectionFormPat"
import { PATList } from "./PatList"

export function FindAllPAT() {
  const { data: pats, isLoading } = useFindAllPat()
  const { mutate: createPat, isPending: creating } = useCreatePat()
  const { mutate: deletePat } = useDeletePat()

  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState("")
  const [expiresAt, setExpiresAt] = useState("")
  const [newToken, setNewToken] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  const handleCreate = () => {
    if (!name) return

    createPat(
      {
        data: {
          name,
          expires_at: expiresAt ? new Date(expiresAt).toISOString() : null,
        },
      },
      {
        onSuccess: (res) => {
          setNewToken(res.token) 
          setShowForm(false)
          setName("")
          setExpiresAt("")
        },
      }
    )
  }

  const handleCopy = () => {
    if (!newToken) return
    navigator.clipboard.writeText(newToken)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="max-w-2xl px-6 py-8">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-sm font-medium text-sidebar-text-active">
            Personal Access Tokens
          </h2>
          <p className="text-[11px] mt-0.5 text-sidebar-text">
            Use PATs to access your workspaces via API or CLI
          </p>
        </div>

        <button
          onClick={() => setShowForm(v => !v)}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded text-[11px] border"
        >
          <Plus size={12} />
          New Token
        </button>
      </div>

      {/* Result */}
      {newToken && (
        <SectionResultPat
          copied={copied}
          handleCopy={handleCopy}
          newToken={newToken}
          setNewToken={setNewToken}
        />
      )}

      {/* Form */}
      {showForm && (
        <SectionFormPat
          name={name}
          setName={setName}
          expiresAt={expiresAt}
          setExpiresAt={setExpiresAt}
          onCreate={handleCreate}
          onCancel={() => setShowForm(false)}
          creating={creating}
        />
      )}

      {/* List */}
      <PATList
        pats={pats}
        isLoading={isLoading}
        onDelete={deletePat}
      />
    </div>
  )
}