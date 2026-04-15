import { Plus, Trash2 } from "lucide-react"
import { styles as s } from "./styles"
import { emptyFile } from "./types"
import { useTemplateDetails } from "./TemplateDetailsContext";

export function TemplateFiles() {
  const {
    filesData,
    showFileForm, setShowFileForm,
    newFile, setNewFile,
    submitFile,
    deleteFile,
  } = useTemplateDetails()

  return (
    <section className="px-5 py-4">
      <div className="flex items-center justify-between mb-3">
        <span className={s.sectionLabel}>files</span>
        <button
          type="button"
          className={s.addButton}
          onClick={() => setShowFileForm(v => !v)}
        >
          <Plus size={9} />
          add file
        </button>
      </div>

      {showFileForm && (
        <div className={s.inlineForm}>
          <div className="grid gap-2" style={{ gridTemplateColumns: "1fr 80px" }}>
            <input
              className={s.inp}
              placeholder="filename (e.g. .bashrc)"
              value={newFile.filename}
              onChange={e => setNewFile(p => ({ ...p, filename: e.target.value }))}
            />
            <input
              className={s.inp}
              placeholder="order"
              type="number"
              value={newFile.sort_order}
              onChange={e => setNewFile(p => ({ ...p, sort_order: Number(e.target.value) }))}
            />
          </div>

          <textarea
            className={s.inp}
            placeholder="file content..."
            rows={4}
            style={{ resize: "vertical" }}
            value={newFile.content}
            onChange={e => setNewFile(p => ({ ...p, content: e.target.value }))}
          />

          <div className="flex justify-end gap-2">
            <button
              type="button"
              className={s.btnCancel}
              onClick={() => {
                setShowFileForm(false)
                setNewFile(emptyFile)
              }}
            >
              cancel
            </button>
            <button
              type="button"
              className={s.btnSave}
              onClick={submitFile}
            >
              save
            </button>
          </div>
        </div>
      )}

      {filesData && filesData.length > 0 && (
        <div className={`${s.tableWrap} mt-2`}>
          <div
            className="grid px-3 py-1.5 border-b border-[#111] bg-[#0a0a0a]"
            style={{ gridTemplateColumns: "1fr 60px 28px" }}
          >
            <span className={s.colHead}>filename</span>
            <span className={s.colHead}>order</span>
            <span />
          </div>

          {filesData.map((f, i) => (
            <div
              key={f.id}
              className={`grid px-3 py-2 items-center ${
                i < filesData.length - 1 ? "border-b border-[#0f0f0f]" : ""
              }`}
              style={{ gridTemplateColumns: "1fr 60px 28px" }}
            >
              <span className="text-[11px] text-[#71717a] font-mono">
                {f.filename}
              </span>

              <span className="text-[10px] text-sidebar-text-muted font-mono">
                #{f.sort_order}
              </span>

              <button
                type="button"
                className={s.dangerBtn}
                onClick={() => deleteFile(f.id)}
                title="Delete file"
              >
                <Trash2 size={11} />
              </button>
            </div>
          ))}
        </div>
      )}
    </section>
  )
}