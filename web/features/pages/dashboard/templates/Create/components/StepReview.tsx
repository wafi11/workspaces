import { CreateTemplateRequest } from "@/types";
import { ReviewRow } from "../../helpers";
import { Badge } from "@/components/ui/badge";

export function StepReview({ values }: { values: CreateTemplateRequest }) {
  return (
    <div className="space-y-4 text-sm">
      <ReviewRow label="Name" value={values.name} />
      <ReviewRow label="Category" value={values.category} />
      <ReviewRow label="Description" value={values.description || "—"} />
      <ReviewRow label="Image" value={values.image || "—"} mono />
      <ReviewRow
        label="Visibility"
        value={
          <Badge variant={values.is_public ? "default" : "secondary"}>
            {values.is_public ? "Public" : "Private"}
          </Badge>
        }
      />
      <ReviewRow
        label="Variables"
        value={`${values.variables.length} defined`}
      />
      <ReviewRow label="Addons" value={`${values.addons.length} attached`} />
    </div>
  );
}
