import { useCreateTemplate } from "@/features/hooks/templates/useCreateTemplate";
import { EmptyState, Field } from "../../helpers";
import { Input } from "@/components/ui/input";
import { Controller } from "react-hook-form";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Plus, Trash2 } from "lucide-react";

export function StepVariables({
  form,
  variables,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
  variables: ReturnType<typeof useCreateTemplate>["variables"];
}) {
  const { register, control } = form;
  return (
    <div className="space-y-3">
      {variables.fields.length === 0 && (
        <EmptyState label="No environment variables. Click Add to define one." />
      )}
      {variables.fields.map((field, i) => (
        <div
          key={field.id}
          className="rounded-lg border border-border bg-muted/20 p-3 space-y-3"
        >
          <div className="grid grid-cols-2 gap-3">
            <Field label="Key">
              <Input
                {...register(`variables.${i}.key`)}
                placeholder="DATABASE_URL"
                className="font-mono text-xs"
              />
            </Field>
            <Field label="Default Value">
              <Input
                {...register(`variables.${i}.default_value`)}
                placeholder="postgres://..."
                className="font-mono text-xs"
              />
            </Field>
          </div>
          <Field label="Description">
            <Input
              {...register(`variables.${i}.description`)}
              placeholder="What is this for?"
            />
          </Field>
          <div className="flex items-center justify-between">
            <label className="flex items-center gap-2 cursor-pointer">
              <Controller
                control={control}
                name={`variables.${i}.required`}
                render={({ field }) => (
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                )}
              />
              <span className="text-xs text-muted-foreground">Required</span>
            </label>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 text-xs text-destructive hover:text-destructive gap-1"
              onClick={() => variables.remove(i)}
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
        onClick={variables.append}
      >
        <Plus className="h-3.5 w-3.5" /> Add Variable
      </Button>
    </div>
  );
}
