import { useFindTemplateWorkspaceForm, useGetTemplateVariables } from "@/features/api";
import type { WorkspaceRequest } from "@/types";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { ActionBtn } from "../ActionButton";
import { Field } from "../Fields";
import { inputStyle } from "../InputStyle";

interface FormWorkspacesProps {
  onSuccess: (ws : WorkspaceRequest) => void;
  templateId? : string
}

export function FormWorkspaces({ onSuccess,templateId }: FormWorkspacesProps) {
  const { data: templates, isLoading } = useFindTemplateWorkspaceForm();
  const [selectedTemplateId, setSelectedTemplateId] = useState<string | undefined>(templateId);
  const {data : templateVariables}  = useGetTemplateVariables(selectedTemplateId)

const { register, handleSubmit, control, setValue, formState: { errors, isSubmitting } } =useForm<WorkspaceRequest>({
  defaultValues: {
    name: "",
    description: "",
    template_id: "",
    limit_ram_mb: 1024,
    limit_cpu_cores: 1.0,
    req_ram_mb: 512,
    req_cpu_cores: 0.5,
    env_vars: {},
  },
})

// update keduanya sekaligus
const handleTemplateChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
  const id = e.target.value
  setSelectedTemplateId(id)
  setValue("template_id", id)  // sync ke RHF
}

  const onSubmit = async (data: WorkspaceRequest) => {
    onSuccess?.(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col w-full gap-4">
      {/* Template */}
      <Field label="Template">
        {isLoading ? (
          <span
            style={{
              fontSize: "12px",
              color: "var(--color-sidebar-text-muted)",
            }}
          >
            Loading...
          </span>
        ) : (
          <select
            value={selectedTemplateId}
            onChange={(e) =>handleTemplateChange(e)}
            style={inputStyle}
          >
            <option value="">Select a template</option>
            {templates?.map((t) => (
              <option key={t.id} value={t.id}>
                {t.name}
              </option>
            ))}
          </select>
        )}
      </Field>

      {/* Name */}
      <Field label="Name" error={errors.name?.message}>
        <input
          {...register("name", { required: "Name is required" })}
          placeholder="my-workspace"
          style={inputStyle}
        />
      </Field>


      {/* Description */}
      <Field label="Description">
        <input
          {...register("description")}
          placeholder="Optional"
          style={inputStyle}
        />
      </Field>

      <Field label="Password" error={errors.password?.message}>
      <input
          {...register("password", { required: "Password is required" })}
          placeholder="******"
          style={inputStyle}
      />
      </Field>

      <Field label="RAM Limit (MB)">
  <input
    type="number"
    {...register("limit_ram_mb", { valueAsNumber: true })}
    style={inputStyle}
  />
</Field>

<Field label="CPU Limit (cores)">
  <input
    type="number"
    step="0.1"
    {...register("limit_cpu_cores", { valueAsNumber: true })}
    style={inputStyle}
  />
</Field>

<Field label="RAM Request (MB)">
  <input
    type="number"
    {...register("req_ram_mb", { valueAsNumber: true })}
    style={inputStyle}
  />
</Field>

<Field label="CPU Request (cores)">
  <input
    type="number"
    step="0.1"
    {...register("req_cpu_cores", { valueAsNumber: true })}
    style={inputStyle}
  />
</Field>

      {/* Template Variables */}
      {templateVariables && templateVariables.length > 0 && (
        <>
          <div
            className="text-[11px] pt-1"
            style={{ color: "var(--color-sidebar-text-muted)" }}
          >
            Template Variables
          </div>
          {templateVariables.map((v) => (
            <Field
              key={v.key}
              label={`${v.key}${v.required ? " *" : ""}`}
            >
              <Controller
                control={control}
                name={`env_vars.${v.key}` as any}
                rules={{
                  required: v.required ? `${v.key} is required` : false,
                }}
                render={({ field }) => (
                  <input
                    {...field}
                    placeholder={v.default_value || ""}
                    style={inputStyle}
                  />
                )}
              />
            </Field>
          ))}
        </>
      )}

      {/* Actions */}
      <div className="flex gap-2 justify-end pt-2">
        <ActionBtn
          label={isSubmitting ? "Creating..." : "Create"}
          variant="default"
          type="submit"
          disabled={isSubmitting || !selectedTemplateId}
        />
      </div>
    </form>
  );
}
