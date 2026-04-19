type Props = {
  name: string
  setName: (v: string) => void
  expiresAt: string
  setExpiresAt: (v: string) => void
  onCreate: () => void
  onCancel: () => void
  creating: boolean
}

export function SectionFormPat({
  name,
  setName,
  expiresAt,
  setExpiresAt,
  onCreate,
  onCancel,
  creating,
}: Props) {
  return (
    <div className="mb-4 p-4 rounded border">
      <div className="flex flex-col gap-3">
        <div>
          <label className="text-[10px] uppercase mb-1">
            Token name
          </label>
          <input
            value={name}
            onChange={e => setName(e.target.value)}
            placeholder="e.g. my-cli-token"
            className="w-full px-3 py-2 rounded text-[12px] border"
          />
        </div>

        <div>
          <label className="text-[10px] uppercase mb-1">
            Expires at (optional)
          </label>
          <input
            type="date"
            value={expiresAt}
            onChange={e => setExpiresAt(e.target.value)}
            className="w-full px-3 py-2 rounded text-[12px] border"
          />
        </div>

        <div className="flex gap-2">
          <button
            onClick={onCreate}
            disabled={!name || creating}
            className="px-3 py-1.5 rounded text-[11px] bg-accent disabled:opacity-40"
          >
            {creating ? "Generating..." : "Generate token"}
          </button>

          <button
            onClick={onCancel}
            className="px-3 py-1.5 rounded text-[11px] border"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  )
}