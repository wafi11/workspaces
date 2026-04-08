"use client"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useTemplateForm } from "@/features/services/templates"
import { useCreateAddonWorkspaces, useGetListWorkspaceForm } from "@/features/services/workspaces/api"
import { useForm, Controller } from "react-hook-form"

interface TemplateCreateProps {
  id: string
}

type FormValues = {
  variables: Record<string, string>
  addons: string
  workspaceId: string
}

export function TemplateCreate({ id }: TemplateCreateProps) {
  const { data: template, isLoading: isTemplateLoading } = useTemplateForm(id)
  const { data: workspaces, isLoading: isWorkspaceLoading } = useGetListWorkspaceForm()
  const { mutate, isPending } = useCreateAddonWorkspaces()

  const { register, handleSubmit, control } = useForm<FormValues>({
    defaultValues: {
      variables: {},
      addons: "",
      workspaceId: "",
    },
  })

  const onSubmit = (values: FormValues) => {
    if (!values.workspaceId) return

    mutate({
      config: Object.entries(values.variables).map(([key, value]) => ({ key, value })),
      workspace_id: values.workspaceId,
      template_addon_id: values.addons,
    })
  }

  if (isTemplateLoading || isWorkspaceLoading) {
    return <div className="p-4 text-center text-muted-foreground animate-pulse">Loading form configuration...</div>
  }

  if (!template) return null

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
      {/* SECTION: HEADER INFO */}
      <div className="pb-4 border-b">
        <Label className="text-[10px] text-muted-foreground uppercase tracking-widest font-bold">Selected Template</Label>
        <h3 className="text-lg font-semibold tracking-tight">{template.template_name}</h3>
      </div>

      <div className="grid gap-6">
        {/* SECTION: WORKSPACE SELECTION */}
        <div className="space-y-2">
          <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Target Workspace</Label>
         <Controller
              control={control}
              name="workspaceId"
              rules={{ required: true }}
              render={({ field }) => (
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger className="w-full">
                    <SelectValue>
                      {/* Force display name berdasarkan id yang tersimpan */}
                      {field.value 
                        ? workspaces?.find((ws) => ws.id === field.value)?.name 
                        : "Select destination workspace"}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    {workspaces?.map((ws) => (
                      <SelectItem key={ws.id} value={ws.id}>
                        {ws.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
        </div>

        {/* SECTION: ADDONS (VERSIONING) */}
        {template.addons && template.addons.length > 0 && (
          <div className="space-y-2">
            <div className="flex flex-col">
              <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Addon / Version</Label>
              <span className="text-[11px] text-muted-foreground mb-1">Select an optional configuration version</span>
            </div>
          <Controller
            control={control}
            name="addons"
            render={({ field }) => (
              <Select onValueChange={field.onChange} value={field.value}>
                <SelectTrigger>
                  <SelectValue>
                    {field.value
                      ? template.addons.find((a) => a.id === field.value)?.name
                      : "Select Versions"}
                  </SelectValue>
                </SelectTrigger>
                <SelectContent>
                  {template.addons.map((addon) => (
                    <SelectItem key={addon.id} value={addon.id}>
                      {addon.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          </div>
        )}
      </div>

      {/* SECTION: DYNAMIC VARIABLES */}
      {template.variables && template.variables.length > 0 && (
        <div className="space-y-4 pt-4 border-t">
          <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Configuration Variables</Label>
          <div className="grid gap-4">
            {template.variables.map((variable) => (
              <div key={variable.key} className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor={variable.key} className="text-sm font-medium">
                    {variable.key.replaceAll("_", " ")}
                  </Label>
                  {variable.required && (
                    <Badge variant="secondary" className="text-[9px] h-4 uppercase bg-destructive/10 text-destructive border-none">
                      Required
                    </Badge>
                  )}
                </div>
                <Input
                  id={variable.key}
                  className="bg-background"
                  placeholder={`Enter value for ${variable.key.toLowerCase()}...`}
                  {...register(`variables.${variable.key}`, {
                    required: variable.required,
                  })}
                />
              </div>
            ))}
          </div>
        </div>
      )}

      {/* ACTION BUTTON */}
      <Button 
        type="submit" 
        className="w-full font-semibold transition-all active:scale-[0.98]" 
        disabled={isPending}
      >
        {isPending ? "Creating..." : "Confirm & Create Workspace"}
      </Button>
    </form>
  )
}