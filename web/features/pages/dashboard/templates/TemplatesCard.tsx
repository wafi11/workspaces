import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { Templates } from "@/types";
import { Box, ChevronRight, Globe, Lock, Pencil } from "lucide-react";
import Link from "next/link";

export function TemplateCard({
  item,
  isActive,
  onClick,
}: {
  item: Templates;
  isActive: boolean;
  onClick: () => void;
}) {
  return (
    <Link
      href={`/dashboard/templates/${item.id}`}
      className={cn(
        "group relative flex flex-col gap-3 rounded-xl border bg-card p-4 cursor-pointer",
        "transition-all duration-200 hover:shadow-md hover:-translate-y-0.5",
        isActive
          ? "border-primary shadow-md ring-1 ring-primary/30"
          : "border-border hover:border-primary/40"
      )}
    >
      {/* Header */}
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-3">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-lg text-xl">
            {item.icon ? (
              <img
                width={100}
                height={100}
                src={item.icon}
                alt={item.name}
                className="size-10 object-cover rounded-lg"
              />
            ) : (
              <Box className="h-5 w-5 text-muted-foreground" />
            )}
          </div>

          <div className="min-w-0">
            <p className="truncate font-semibold text-sm text-foreground leading-tight">
              {item.name}
            </p>
            <p className="text-xs text-muted-foreground mt-0.5">
              {item.category}
            </p>
          </div>
        </div>

        <ChevronRight
          className={cn(
            "h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200",
            "group-hover:translate-x-0.5 group-hover:text-primary"
          )}
        />
      </div>

      {/* Description */}
      <p className="line-clamp-2 text-xs text-muted-foreground leading-relaxed">
        {item.description || "No description provided."}
      </p>

      {/* Footer */}
      <div className="flex items-center gap-2 mt-auto pt-1 relative">
        <Badge
          variant={item.is_public ? "default" : "secondary"}
          className="gap-1 text-xs px-2 py-0.5"
        >
          {item.is_public ? (
            <Globe className="h-3 w-3" />
          ) : (
            <Lock className="h-3 w-3" />
          )}
          {item.is_public ? "Public" : "Private"}
        </Badge>
        <button
          onClick={onClick}
          className="ml-auto absolute z-20 bottom-1 right-2 flex items-center gap-1 text-xs text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity hover:text-foreground"
        >
          <Pencil className="h-3 w-3" />
          Edit
        </button>
      </div>
    </Link>
  );
}
