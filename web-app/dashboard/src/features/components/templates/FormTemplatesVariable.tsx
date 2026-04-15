import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import type { useCreateTemplate } from "@/hooks/useCreateTemplates";
import { Field } from "../Fields";
import { inputStyle } from "../InputStyle";
import { Trash2, Plus } from "lucide-react";

export function FormTemplatesVariables({
  form,
  variables,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
  variables: ReturnType<typeof useCreateTemplate>["variables"];
}) {
  const { register, formState: { errors } } = form;
  const { fields, append, remove } = variables;

  return (
    <>
      <TopbarAdmin title="Variables" classNameFont="text-md font-medium" />
      <div className="flex flex-col gap-3 px-4 py-4 overflow-y-auto">
        {fields.map((field, index) => (
          <div
            key={field.id}
            className="flex flex-col gap-2 p-3 rounded-md"
            style={{ border: "1px solid var(--color-sidebar-border)" }}
          >
            <div className="grid grid-cols-2 gap-2">
              <Field label="Key" error={errors.variables?.[index]?.key?.message}>
                <input
                  {...register(`variables.${index}.key`, { required: "Required" })}
                  placeholder="DATABASE_URL"
                  style={inputStyle}
                />
              </Field>
              <Field label="Default Value">
                <input
                  {...register(`variables.${index}.default_value`)}
                  placeholder="postgres://..."
                  style={inputStyle}
                />
              </Field>
            </div>
            <Field label="Description">
              <input
                {...register(`variables.${index}.description`)}
                placeholder="description"
                style={inputStyle}
              />
            </Field>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id={`required-${index}`}
                  {...register(`variables.${index}.required`)}
                />
                <label
                  htmlFor={`required-${index}`}
                  style={{ fontSize: 12, color: "var(--color-sidebar-text-muted)" }}
                >
                  Required
                </label>
              </div>
              <button
                type="button"
                onClick={() => remove(index)}
                style={{ color: "#f87171", background: "none", border: "none", cursor: "pointer", padding: 4 }}
              >
                <Trash2 size={14} />
              </button>
            </div>
          </div>
        ))}

        <button
          type="button"
          onClick={append}
          className="flex items-center gap-2 w-full justify-center py-2 rounded-md"
          style={{
            border: "1px dashed var(--color-sidebar-border)",
            color: "var(--color-sidebar-text-muted)",
            background: "none",
            cursor: "pointer",
            fontSize: 13,
          }}
        >
          <Plus size={14} />
          Add variable
        </button>
      </div>
    </>
  );
}