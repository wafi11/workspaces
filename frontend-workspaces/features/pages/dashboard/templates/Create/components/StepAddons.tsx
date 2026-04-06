import { useCreateTemplate } from "@/features/hooks/templates/useCreateTemplate";
import { EmptyState, Field } from "../../helpers";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Plus, Trash2 } from "lucide-react";

export function StepAddons({
  form,
  addons,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
  addons: ReturnType<typeof useCreateTemplate>["addons"];
}) {
  const { register } = form;
  return (
    <div className="space-y-3">
      {addons.fields.length === 0 && (
        <EmptyState label="No addons. Click Add to attach a sidecar service." />
      )}
      {addons.fields.map((field, i) => (
        <div
          key={field.id}
          className="rounded-lg border border-border bg-muted/20 p-3 space-y-3"
        >
          <div className="grid grid-cols-2 gap-3">
            <Field label="Name">
              <Input {...register(`addons.${i}.name`)} placeholder="postgres" />
            </Field>
            <Field label="Image">
              <Input
                {...register(`addons.${i}.image`)}
                placeholder="postgres:15-alpine"
                className="font-mono text-xs"
              />
            </Field>
          </div>
          <Field label="Description">
            <Input
              {...register(`addons.${i}.description`)}
              placeholder="Sidecar database"
            />
          </Field>
          <div className="flex justify-end">
            <Button
              variant="ghost"
              size="sm"
              className="h-7 text-xs text-destructive hover:text-destructive gap-1"
              onClick={() => addons.remove(i)}
            >
              <Trash2 className="h-3 w-3" /> Remove
            </Button>
          </div>
        </div>
      ))}
      <Button
        variant="outline"
        size="sm"
        className="w-full gap-1"
        onClick={addons.append}
      >
        <Plus className="h-3.5 w-3.5" /> Add Addon
      </Button>
    </div>
  );
}
