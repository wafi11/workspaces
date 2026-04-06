import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { Templates } from "@/types";
import { X } from "lucide-react";
import { Field } from "./helpers";
import { useUpdateTemplates } from "@/features/services/templates/api";

export function EditForm({
  form,
  onChange,
  onClose,
}: {
  form: Partial<Templates>;
  onChange: (key: keyof Templates, value: string | boolean) => void;
  onClose: () => void;
}) {
  const { mutate } = useUpdateTemplates(form.id as string);

  const handleSubmit = () => {
    mutate(form);
  };
  return (
    <>
      {/* Header */}
      <div className="flex items-center justify-between border-b border-border px-6 py-4">
        <div>
          <h2 className="font-semibold text-foreground">Edit Template</h2>
          <p className="text-xs text-muted-foreground mt-0.5">{form.name}</p>
        </div>
        <Button variant="ghost" size="icon" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      {/* Body */}
      <div className="flex-1 overflow-y-auto px-6 py-5 space-y-5">
        <Field label="Name">
          <Input
            value={form.name ?? ""}
            onChange={(e) => onChange("name", e.target.value)}
            placeholder="Template name"
          />
        </Field>

        <Field label="Category">
          <Input
            value={form.category ?? ""}
            onChange={(e) => onChange("category", e.target.value)}
            placeholder="e.g. Database, IDE, DevOps"
          />
        </Field>

        <Field label="Description">
          <Textarea
            value={form.description ?? ""}
            onChange={(e) => onChange("description", e.target.value)}
            placeholder="Short description of the template"
            rows={3}
          />
        </Field>

        <Field label="Template URL">
          <Input
            value={form.template_url ?? ""}
            onChange={(e) => onChange("template_url", e.target.value)}
            placeholder="https://..."
          />
        </Field>

        <Field label="Docker Image">
          <Input
            value={form.image ?? ""}
            onChange={(e) => onChange("image", e.target.value)}
            placeholder="bitnami/postgres:2.1"
          />
        </Field>

        <Field label="Icon Url">
          <Input
            value={form.icon ?? ""}
            onChange={(e) => onChange("icon", e.target.value)}
            placeholder="e.g. 🐘"
          />
        </Field>

        <div className="flex items-center justify-between rounded-lg border border-border bg-muted/30 px-4 py-3">
          <div>
            <p className="text-sm font-medium text-foreground">Public</p>
            <p className="text-xs text-muted-foreground">
              Visible to all users
            </p>
          </div>
          <Switch
            checked={form.is_public ?? false}
            onCheckedChange={(val) => onChange("is_public", val)}
          />
        </div>
      </div>

      {/* Footer */}
      <div className="flex gap-2 border-t border-border px-6 py-4">
        <Button variant="outline" className="flex-1" onClick={onClose}>
          Cancel
        </Button>
        <Button onClick={handleSubmit} className="flex-1">
          Save Changes
        </Button>
      </div>
    </>
  );
}
