import { useFindTemplateWorkspaceForm, useTemplateForm } from "@/features/api";
import type { WorkspaceRequest } from "@/types";
import * as Dialog from "@radix-ui/react-dialog";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { ActionBtn } from "../ActionButton";
import { Field } from "../Fields";
import { inputStyle } from "../InputStyle";

interface FormWorkspacesProps {
  onSuccess: () => void;
}

export function FormWorkspaces({ onSuccess }: FormWorkspacesProps) {
  const { data: templates, isLoading } = useFindTemplateWorkspaceForm();
  const [selectedTemplateId, setSelectedTemplateId] = useState("");
  const { data: templateDetail } = useTemplateForm(selectedTemplateId);

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<WorkspaceRequest>({
    defaultValues: {
      name: "",
      description: "",
      password: "",
      env_vars: {},
    },
  });

  const onSubmit = async (data: WorkspaceRequest) => {
    console.log({ ...data, template_id: selectedTemplateId });
    // TODO: mutation POST /workspaces
    onSuccess?.();
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
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
            onChange={(e) => setSelectedTemplateId(e.target.value)}
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

      {/* Template Variables */}
      {templateDetail?.variables && templateDetail.variables.length > 0 && (
        <>
          <div
            className="text-[11px] pt-1"
            style={{ color: "var(--color-sidebar-text-muted)" }}
          >
            Template Variables
          </div>
          {/* {templateDetail.variables.map((v) => (
            <Field
              key={v.name}
              label={`${v.display_name || v.name}${v.required ? " *" : ""}`}
            >
              <Controller
                control={control}
                name={`env_vars.${v.name}` as any}
                rules={{
                  required: v.required ? `${v.name} is required` : false,
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
          ))} */}
        </>
      )}

      {/* Actions */}
      <div className="flex gap-2 justify-end pt-2">
        <Dialog.Close asChild>
          <ActionBtn label="Cancel" variant="default" />
        </Dialog.Close>
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
