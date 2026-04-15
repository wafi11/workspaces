import { Plus, Trash2 } from "lucide-react"
import { styles as s } from "./styles"
import { emptyAddon } from "./types"
import { useTemplateDetails } from "./TemplateDetailsContext"

export function SectionAddOn() {
  const {
    addons,
    showAddonForm,
    setShowAddonForm,
    newAddon,
    setNewAddon,
    submitAddon,
    deleteAddon,
  } = useTemplateDetails()

  return (
    <section className="px-5 py-4 border-b border-[#0f0f0f]">
      <div className="flex items-center justify-between mb-3">
        <span className={s.sectionLabel}>addons</span>
        <button
          type="button"
          className={s.addButton}
          onClick={() => setShowAddonForm(v => !v)}
        >
          <Plus size={9} />
          add addon
        </button>
      </div>

      {showAddonForm && (
        <div className={s.inlineForm}>
          <div className="grid grid-cols-2 gap-2">
            <input
              className={s.inp}
              placeholder="name"
              value={newAddon.name}
              onChange={e => setNewAddon(p => ({ ...p, name: e.target.value }))}
            />
            <input
              className={s.inp}
              placeholder="image"
              value={newAddon.image}
              onChange={e => setNewAddon(p => ({ ...p, image: e.target.value }))}
            />
          </div>

          <input
            className={s.inp}
            placeholder="description"
            value={newAddon.description}
            onChange={e => setNewAddon(p => ({ ...p, description: e.target.value }))}
          />

          <div className="flex justify-end gap-2">
            <button
              type="button"
              className={s.btnCancel}
              onClick={() => {
                setShowAddonForm(false)
                setNewAddon(emptyAddon)
              }}
            >
              cancel
            </button>
            <button type="button" className={s.btnSave} onClick={submitAddon}>
              save
            </button>
          </div>
        </div>
      )}

      {addons && addons.length > 0 && (
        <div className="grid grid-cols-2 gap-1.5 mt-2">
          {addons.map(a => (
            <div key={a.id} className="flex justify-between px-3 py-2.5 bg-[#0d0d0d] border border-[#141414] rounded-sm">
              <div>
                <span>{a.name}</span>
              </div>
              <button onClick={() => deleteAddon(a.id)}>
                <Trash2 size={11} />
              </button>
            </div>
          ))}
        </div>
      )}
    </section>
  )
}