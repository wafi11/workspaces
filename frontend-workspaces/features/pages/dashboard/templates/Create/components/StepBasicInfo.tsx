import { useCreateTemplate } from "@/features/hooks/templates/useCreateTemplate";
import { Field } from "../../helpers";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Controller } from "react-hook-form";
import { Switch } from "@/components/ui/switch";

export function StepBasicInfo({
  form,
}: {
  form: ReturnType<typeof useCreateTemplate>["form"];
}) {
  const {
    register,
    control,
    formState: { errors },
  } = form;
  return (
    <div className="space-y-4">
      <Field label="Name" error={errors.name?.message}>
        <Input
          {...register("name", { required: "Name is required" })}
          placeholder="VSCode Server"
        />
      </Field>
      <Field label="Category" error={errors.category?.message}>
        <Input
          {...register("category", { required: "Category is required" })}
          placeholder="IDE, Database, DevOps…"
        />
      </Field>
      <Field label="Description">
        <Textarea
          {...register("description")}
          placeholder="What does this template provide?"
          rows={3}
        />
      </Field>
      <Field label="Docker Image">
        <Input
          {...register("image")}
          placeholder="codercom/code-server:latest"
          className="font-mono text-xs"
        />
      </Field>
      <Field label="Icon URL">
        <Input
          {...register("icon")}
          placeholder="https://example.com/logo.png"
        />
      </Field>
      <div className="flex items-center justify-between rounded-lg border border-border bg-muted/20 px-4 py-3">
        <div>
          <p className="text-sm font-medium">Public</p>
          <p className="text-xs text-muted-foreground">Visible to all users</p>
        </div>
        <Controller
          control={control}
          name="is_public"
          render={({ field }) => (
            <Switch checked={field.value} onCheckedChange={field.onChange} />
          )}
        />
      </div>
    </div>
  );
}
