import { Plus, Pencil, Trash2 } from "lucide-react"
import { styles as s } from "./styles"
import { emptyVar } from "./types"
import { useTemplateDetails } from "./TemplateDetailsContext";

export function TemplateVariable() {
  const {
    variables,
    showVarForm, setShowVarForm,
    newVar, setNewVar,
    editingVarId,
    editVarValue, setEditVarValue,

    startEditVar,
    saveVar,
    cancelEditVar,
    submitVar,
    deleteVar,
  } = useTemplateDetails()

  return (
    <section className="px-5 py-4 border-b border-[#0f0f0f]">
      <div className="flex items-center justify-between mb-3">
        <span className={s.sectionLabel}>variables</span>
        <button
          type="button"
          className={s.addButton}
          onClick={() => setShowVarForm(v => !v)}
        >
          <Plus size={9} />
          add variable
        </button>
      </div>

      {/* Add variable form */}
      {showVarForm && (
        <div className={s.inlineForm}>
          <div className="grid grid-cols-2 gap-2">
            <input
              className={s.inp}
              placeholder="KEY_NAME"
              value={newVar.key}
              onChange={e => setNewVar(p => ({ ...p, key: e.target.value }))}
            />
            <input
              className={s.inp}
              placeholder="default value"
              value={newVar.default_value}
              onChange={e => setNewVar(p => ({ ...p, default_value: e.target.value }))}
            />
          </div>

          <input
            className={s.inp}
            placeholder="description"
            value={newVar.description}
            onChange={e => setNewVar(p => ({ ...p, description: e.target.value }))}
          />

          <div className="flex items-center justify-between">
            <label className="flex items-center gap-2 cursor-pointer text-[10px] text-sidebar-text font-mono select-none">
              <input
                type="checkbox"
                checked={newVar.required}
                onChange={e => setNewVar(p => ({ ...p, required: e.target.checked }))}
                className="accent-[#4ade80]"
              />
              required
            </label>

            <div className="flex gap-2">
              <button
                type="button"
                className={s.btnCancel}
                onClick={() => {
                  setShowVarForm(false)
                  setNewVar(emptyVar)
                }}
              >
                cancel
              </button>
              <button type="button" className={s.btnSave} onClick={submitVar}>
                save
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Variables table */}
      {variables && variables.length > 0 && (
        <div className={`${s.tableWrap} mt-2`}>
          <div
            className="grid px-3 py-1.5 border-b border-[#111] bg-[#0a0a0a]"
            style={{ gridTemplateColumns: "140px 120px 1fr auto 28px 28px" }}
          >
            <span className={s.colHead}>key</span>
            <span className={s.colHead}>default</span>
            <span className={s.colHead}>description</span>
            <span className={s.colHead} />
            <span />
            <span />
          </div>

          {variables.map((v, i) => (
            <div
              key={v.id}
              className={i < variables.length - 1 ? "border-b border-[#0f0f0f]" : ""}
            >
              {editingVarId === v.id ? (
                <div className="flex flex-col gap-2 px-3 py-2.5">
                  <div className="grid grid-cols-2 gap-2">
                    <input
                      className={s.inp}
                      value={editVarValue.key}
                      onChange={e => setEditVarValue(p => ({ ...p, key: e.target.value }))}
                    />
                    <input
                      className={s.inp}
                      value={editVarValue.default_value}
                      onChange={e => setEditVarValue(p => ({ ...p, default_value: e.target.value }))}
                    />
                  </div>

                  <input
                    className={s.inp}
                    value={editVarValue.description}
                    onChange={e => setEditVarValue(p => ({ ...p, description: e.target.value }))}
                  />

                  <div className="flex justify-between">
                    <button onClick={cancelEditVar}>cancel</button>
                    <button onClick={() => saveVar(v.id)}>save</button>
                  </div>
                </div>
              ) : (
                <div
                  className="grid px-3 py-2 items-center"
                  style={{ gridTemplateColumns: "140px 120px 1fr auto 28px 28px" }}
                >
                  <span className="text-sm ">{v.key}</span>
                  <span className="text-xs">{v.default_value || "—"}</span>
                  <span className="text-xs">{v.description || "—"}</span>

                  <button onClick={() => startEditVar(v.id, v)}>
                    <Pencil size={11} />
                  </button>

                  <button onClick={() => deleteVar(v.id)}>
                    <Trash2 size={11} />
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </section>
  )
}