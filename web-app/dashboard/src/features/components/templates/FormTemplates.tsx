import { dataCategoryTemplates } from "@/data/dataCategoryTemplates";
import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import { useCreateTemplate } from "@/hooks/useCreateTemplates";
import { Field } from "../Fields";
import { inputStyle } from "../InputStyle";

export function FormTemplates ({form,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
}) {
    const {
        register,
        formState: { errors },
    } = form 


    return (
       <>
        <TopbarAdmin title="Form Templates" classNameFont="text-md text-md font-medium" isUsedButtonNotification={false}/>
        <div className="flex flex-col gap-4 px-4 py-4 overflow-y-auto">
            {/* Row 1 - 2 kolom */}
            <div className="grid grid-cols-2 gap-3">
                <Field label="Name" error={errors.name?.message}>
                <input {...register("name", { required: "Required" })} placeholder="development" style={inputStyle} />
                </Field>
                <Field label="Icon" error={errors.icon?.message}>
                <input {...register("icon", { required: "Required" })} placeholder="🐳" style={inputStyle} />
                </Field>
            </div>

            {/* Description - full width */}
            <Field label="Description" error={errors.description?.message}>
                <textarea
                {...register("description")}
                placeholder="Describe this template..."
                rows={3}
                style={{ ...inputStyle, resize: 'none' }}
                />
            </Field>

            {/* Row 2 - 2 kolom */}
            <div className="grid grid-cols-2 gap-3">
                <Field label="Category">
                <select {...register("category")} style={inputStyle}>
                    <option value="">Select category</option>
                    {dataCategoryTemplates.map((t) => (
                    <option key={t} value={t}>{t.toUpperCase()}</option>
                    ))}
                </select>
                </Field>
                <Field label="Visibility">
                <div className="flex items-center gap-2 h-[30px]">
                    <input type="checkbox" {...register("is_public")} id="is_public" />
                    <label htmlFor="is_public" style={{ fontSize: 13, color: 'var(--color-sidebar-text-active)' }}>
                    Public template
                    </label>
                </div>
                </Field>
            </div>
            </div>
       </>
    )
}