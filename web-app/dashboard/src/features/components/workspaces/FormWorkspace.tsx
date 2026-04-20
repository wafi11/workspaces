import {
  useFindTemplateWorkspaceForm,
  useGetTemplateVariables,
} from "@/features/api";
import type { WorkspaceRequest } from "@/types";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { ActionBtn } from "../ActionButton";
import { Field } from "../Fields";
import { inputStyle } from "../InputStyle";

interface FormWorkspacesProps {
  onSuccess: (ws: WorkspaceRequest) => void;
  templateId?: string;
}

export function FormWorkspaces({ onSuccess, templateId }: FormWorkspacesProps) {
  const { data: templates, isLoading } = useFindTemplateWorkspaceForm();
  const [selectedTemplateId, setSelectedTemplateId] = useState<
    string | undefined
  >(templateId);
  const { data: templateVariables } =
    useGetTemplateVariables(selectedTemplateId);

  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<WorkspaceRequest>({
    defaultValues: {
      name: "",
      description: "",
      template_id: templateId ?? "",
      type_time_duration: "minutes",
      time_duration: 1,
      limit_ram_mb: 1024,
      limit_cpu_cores: 1.0,
      req_ram_mb: 512,
      req_cpu_cores: 0.5,
      env_vars: {},
    },
  });

  const handleTemplateChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const id = e.target.value;
    setSelectedTemplateId(id);
    setValue("template_id", id);
  };

  const onSubmit = async (data: WorkspaceRequest) => {
    onSuccess?.(data);
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="flex flex-col w-full gap-4"
    >
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
            onChange={handleTemplateChange}
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

      {/* Time Duration + Type */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 w-full justify-between">
        <Field label="Time Duration" error={errors.time_duration?.message}>
          <input
            {...register("time_duration", {
              required: "Time Duration is required",
              valueAsNumber: true,
              min: { value: 1, message: "Minimum 1" },
            })}
            placeholder="1"
            type="number"
            className="w-full"
            style={inputStyle}
          />
        </Field>

        <Field label="Duration Type" error={errors.type_time_duration?.message}>
          <select
            {...register("type_time_duration", { required: true })}
            style={inputStyle}
            className="w-full"
          >
            <option value="minutes">Minutes</option>
            <option value="hours">Hours</option>
            <option value="days">Days</option>
          </select>
        </Field>
      </div>

      {/* Resources */}
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
            <Field key={v.key} label={`${v.key}${v.required ? " *" : ""}`}>
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
