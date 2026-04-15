// ── inline form types ──────────────────────────────────────────
export type NewVariable = { key: string; default_value: string; required: boolean; description: string }
export type NewAddon    = { name: string; image: string; description: string }
export type NewFile     = { filename: string; sort_order: number; content: string }

export const emptyVar:   NewVariable = { key: "", default_value: "", required: false, description: "" }
export const emptyAddon: NewAddon    = { name: "", image: "", description: "" }
export const emptyFile:  NewFile     = { filename: "", sort_order: 1, content: "" }
