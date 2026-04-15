import { useState } from "react"
import {
  useCreateTemplateAddOn,
  useCreateTemplateVariables,
  useDeleteTemplateAddOn,
  useDeleteTemplateVariables,
  useGetTemplateAddOns,
  useGetTemplateVariables,
  useTemplateDetails,
  useUpdateTemplateVariables,
} from "@/features/api"
import { useCreateTemplateFiles, useDeleteTemplateFiles, useGetTemplateFiles } from "@/features/api/template-files"
import type { TemplateEditVariable } from "@/types"
import { emptyAddon, emptyFile, emptyVar, type NewAddon, type NewFile, type NewVariable } from "./types"

export function useTemplateDetailsHook(templateId: string) {
  // ── Data fetching ─────────────────────────────────────────
  const { data }            = useTemplateDetails(templateId)
  const { data: addons }    = useGetTemplateAddOns(templateId)
  const { data: variables } = useGetTemplateVariables(templateId)
  const { data: filesData } = useGetTemplateFiles(templateId)

  // ── Mutations ─────────────────────────────────────────────
  const { mutate: deleteAddon }    = useDeleteTemplateAddOn(templateId)
  const { mutate: deleteVariable } = useDeleteTemplateVariables(templateId)
  const { mutate: updateVariable } = useUpdateTemplateVariables(templateId)
  const { mutate: createVariable } = useCreateTemplateVariables(templateId)
  const { mutate: createAddon }    = useCreateTemplateAddOn(templateId)
  const { mutate: createFile }     = useCreateTemplateFiles(templateId)
  const { mutate: deleteFile }     = useDeleteTemplateFiles(templateId)

  // ── Edit variable state ───────────────────────────────────
  const [editingVarId, setEditingVarId] = useState<string | null>(null)
  const [editVarValue, setEditVarValue] = useState<TemplateEditVariable>({
    key: "",
    default_value: "",
    required: true,
    description: "",
  })

  // ── Inline form visibility ────────────────────────────────
  const [showVarForm,   setShowVarForm]   = useState(false)
  const [showAddonForm, setShowAddonForm] = useState(false)
  const [showFileForm,  setShowFileForm]  = useState(false)

  // ── New entity drafts ─────────────────────────────────────
  const [newVar,   setNewVar]   = useState<NewVariable>(emptyVar)
  const [newAddon, setNewAddon] = useState<NewAddon>(emptyAddon)
  const [newFile,  setNewFile]  = useState<NewFile>(emptyFile)

  // ── Variable handlers ─────────────────────────────────────
  const startEditVar = (id: string, v: TemplateEditVariable) => {
    setEditingVarId(id)
    setEditVarValue({
      key: v.key,
      required: v.required,
      default_value: v.default_value ?? "",
      description: v.description ?? "",
    })
  }

  const saveVar = (id: string) =>
    updateVariable({ id, req: editVarValue }, { onSuccess: () => setEditingVarId(null) })

  const cancelEditVar = () => setEditingVarId(null)

  const submitVar = () =>
    createVariable(newVar, {
      onSuccess: () => { setNewVar(emptyVar); setShowVarForm(false) },
    })

  const deleteVar = (id: string) => deleteVariable(id)

  // ── Addon handlers ────────────────────────────────────────
  const submitAddon = () =>
    createAddon(newAddon, {
      onSuccess: () => { setNewAddon(emptyAddon); setShowAddonForm(false) },
    })

  // ── File handlers ─────────────────────────────────────────
  const submitFile = () =>
    createFile(newFile, {
      onSuccess: () => { setNewFile(emptyFile); setShowFileForm(false) },
    })

  return {
    // Data
    data,
    addons,
    variables,
    filesData,

    // Variable form
    showVarForm,
    setShowVarForm,
    newVar,
    setNewVar,
    editingVarId,
    editVarValue,
    setEditVarValue,

    // Addon form
    showAddonForm,
    setShowAddonForm,
    newAddon,
    setNewAddon,

    // File form
    showFileForm,
    setShowFileForm,
    newFile,
    setNewFile,

    // Handlers
    startEditVar,
    saveVar,
    cancelEditVar,
    submitVar,
    deleteVar,
    deleteAddon,
    submitAddon,
    submitFile,
    deleteFile,
  }
}