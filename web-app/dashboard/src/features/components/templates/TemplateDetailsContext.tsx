import { createContext, useContext } from "react"
import { useTemplateDetailsHook } from "./hooks"

type TemplateDetailsContextType = ReturnType<typeof useTemplateDetailsHook>

const TemplateDetailsContext = createContext<TemplateDetailsContextType | null>(null)

export function TemplateDetailsProvider({
  templateId,
  children,
}: {
  templateId: string
  children: React.ReactNode
}) {
  const value = useTemplateDetailsHook(templateId)

  return (
    <TemplateDetailsContext.Provider value={value}>
      {children}
    </TemplateDetailsContext.Provider>
  )
}

export function useTemplateDetails() {
  const ctx = useContext(TemplateDetailsContext)
  if (!ctx) {
    throw new Error("useTemplateDetails must be used inside TemplateDetailsProvider")
  }
  return ctx
}