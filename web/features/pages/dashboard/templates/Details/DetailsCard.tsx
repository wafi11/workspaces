import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Templates } from "@/types";
import { formatDate } from "@/utils";
import { Copy, Database, Globe, Tag } from "lucide-react";

interface TemplateDetailsProps {
  template: Templates;
}

export function TemplateDetailsCard({ template }: TemplateDetailsProps) {
  return (
    <div className="rounded-lg border bg-card p-6 space-y-5">
      {/* Header */}
      <div className="flex gap-4 items-start">
        <div className="w-14 h-14 rounded-lg border bg-muted flex items-center justify-center shrink-0 overflow-hidden">
          {template.icon ? (
            <img
              src={template.icon}
              alt={template.name}
              width={40}
              height={40}
              className="object-contain"
            />
          ) : (
            <Database className="w-6 h-6 text-muted-foreground" />
          )}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <h2 className="text-base font-semibold">{template.name}</h2>
            <Badge variant={template.is_public ? "default" : "secondary"}>
              {template.is_public ? "Public" : "Private"}
            </Badge>
            <Badge variant="outline">{template.category}</Badge>
          </div>
          {template.description && (
            <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
              {template.description}
            </p>
          )}
        </div>

        <p className="text-sm">{formatDate(template.created_at)}</p>
      </div>

      {/* Divider */}
      <div className="h-px bg-border" />

      {/* Meta rows */}
      <div className="space-y-3">
        {/* Image */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground w-24 shrink-0">
            <Tag className="w-3.5 h-3.5" />
            Image
          </div>
          <p className="text-sm font-mono truncate">{template.image}</p>
        </div>

        {/* Template URL */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5 text-xs text-muted-foreground w-24 shrink-0">
            <Globe className="w-3.5 h-3.5" />
            URL
          </div>
          <div className="flex items-center gap-1.5 flex-1 min-w-0">
            <p className="text-sm font-mono truncate flex-1">
              {template.template_url}
            </p>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 shrink-0"
              onClick={() =>
                navigator.clipboard.writeText(template.template_url)
              }
            >
              <Copy className="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
