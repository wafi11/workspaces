import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import type { useCreateTemplate } from "@/hooks/useCreateTemplates";
import { Field } from "../Fields";
import { inputStyle } from "../InputStyle";
import { Trash2, Plus } from "lucide-react";

export function FormTemplatesFiles({
  form,
  files,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
  files: ReturnType<typeof useCreateTemplate>["files"];
}) {
  const { register, formState: { errors } } = form;
  const { fields, append, remove } = files;

  return (
    <>
      <TopbarAdmin title="Files" classNameFont="text-md font-medium" isUsedButtonNotification={false}/>
      <div className="flex flex-col gap-3 px-4 py-4 overflow-y-auto">
        {fields.map((field, index) => (
          <div
            key={field.id}
            className="flex items-center gap-2 p-3 rounded-md"
            style={{ border: "1px solid var(--color-sidebar-border)" }}
          >
            <div className="flex-1">
              <Field label="Filename" error={errors.files?.[index]?.filename?.message}>
                <input
                  {...register(`files.${index}.filename`, { required: "Required" })}
                  placeholder="docker-compose.yml"
                  style={inputStyle}
                />
              </Field>
            </div>
            <div style={{ width: 80 }}>
              <Field label="Order">
                <input
                  type="number"
                  {...register(`files.${index}.sort_order`, { valueAsNumber: true })}
                  placeholder="0"
                  style={inputStyle}
                />
              </Field>
            </div>
            <button
              type="button"
              onClick={() => remove(index)}
              style={{ color: "#f87171", background: "none", border: "none", cursor: "pointer", padding: 4, marginTop: 14 }}
            >
              <Trash2 size={14} />
            </button>
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
          Add file
        </button>
      </div>
    </>
  );
}