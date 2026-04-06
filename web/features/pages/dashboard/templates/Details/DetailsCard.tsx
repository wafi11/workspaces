import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Templates } from "@/types";
import { formatDate } from "@/utils/relativeTime";
import { Calendar, Copy, Database, Globe, Tag } from "lucide-react";

interface TemplateDetailsProps {
  template: Templates;
}
export function TemplateDetailsCard({ template }: TemplateDetailsProps) {
  return (
    <div className="p-6 space-y-6">
      {/* Hero card */}
      <div className="rounded-lg border bg-card p-6 flex gap-5 items-start">
        <div className="w-16 h-16 rounded-lg border bg-muted flex items-center justify-center shrink-0 overflow-hidden">
          {template.icon ? (
            <img
              src={template.icon}
              alt={template.name}
              width={48}
              height={48}
              className="object-contain"
            />
          ) : (
            <Database className="w-7 h-7 text-muted-foreground" />
          )}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <h2 className="text-lg font-semibold">{template.name}</h2>
            <Badge variant={template.is_public ? "default" : "secondary"}>
              {template.is_public ? "Public" : "Private"}
            </Badge>
            <Badge variant="outline">{template.category}</Badge>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            {template.description}
          </p>
        </div>
      </div>

      {/* Detail grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="rounded-lg border bg-card p-4 space-y-1">
          <p className="text-xs text-muted-foreground flex items-center gap-1.5">
            <Tag className="w-3.5 h-3.5" /> Image
          </p>
          <p className="text-sm font-mono truncate">{template.image}</p>
        </div>

        <div className="rounded-lg border bg-card p-4 space-y-1">
          <p className="text-xs text-muted-foreground flex items-center gap-1.5">
            <Calendar className="w-3.5 h-3.5" /> Created at
          </p>
          <p className="text-sm">
            {formatDate(template.created_at)}{" "}
            {/* Assuming created_at is ISO string */}
          </p>
        </div>

        <div className="rounded-lg border bg-card p-4 space-y-1 sm:col-span-2">
          <p className="text-xs text-muted-foreground flex items-center gap-1.5">
            <Globe className="w-3.5 h-3.5" /> Template URL
          </p>
          <div className="flex items-center gap-2">
            <p className="text-sm font-mono truncate flex-1">
              {template.template_url}
            </p>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 shrink-0"
              onClick={() =>
                navigator.clipboard.writeText(template.template_url)
              }
            >
              <Copy className="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      </div>

      <Separator />
    </div>
  );
}
