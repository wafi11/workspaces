import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useCreateTemplate } from "@/features/hooks/templates/useCreateTemplate";
import { Plus } from "lucide-react";
import { EmptyState, Field } from "../../helpers";

export function StepFiles({
  form,
  files,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
  files: ReturnType<typeof useCreateTemplate>["files"];
}) {
  const { register, control } = form;
  return (
    <div className="space-y-3">
      {files.fields.length === 0 && (
        <EmptyState label="No environment files. Click Add to define one." />
      )}
      {files.fields.map((field, i) => (
        <div
          key={field.id}
          className="rounded-lg border border-border bg-muted/20 p-3 space-y-3"
        >
          <div className="grid grid-cols-2 gap-3">
            <Field label="Filename">
              <Input
                {...register(`files.${i}.filename`)}
                placeholder="terminal.yml"
                className="font-mono text-xs"
              />
            </Field>
            <Field label="Sort Order">
              <Input
                {...register(`files.${i}.sort_order`,{
                    valueAsNumber: true,
                })}
                placeholder="1"
                type="number"
                className="font-mono text-xs"
              />
            </Field>
          </div>
        </div>
      ))}
      <Button
        variant="outline"
        size="sm"
        className="w-full gap-1"
        onClick={files.append}
      >
        <Plus className="h-3.5 w-3.5" /> Add Variable
      </Button>
    </div>
  );
}
