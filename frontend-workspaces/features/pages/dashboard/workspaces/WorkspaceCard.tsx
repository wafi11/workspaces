import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Workspaces } from "@/types/workspaces";
import { useRouter } from "next/navigation";
import { statusBadge } from "./Status";

export function WorkspaceCard({ workspace: w }: { workspace: Workspaces }) {
  const badge = statusBadge(w.status);
  const router = useRouter();
  return (
    <div className="relative rounded-xl border border-border bg-background p-5 hover:border-border/60 transition-colors overflow-hidden">
      <div className="pl-2">
        <div className="flex items-start justify-between mb-3">
          <div>
            <p className="font-medium text-sm text-foreground">{w.name}</p>
            <p className="font-mono text-[11px] text-muted-foreground mt-0.5">
              {w.template_id}
            </p>
          </div>
          <span
            className={`font-mono text-[11px] font-medium px-2 py-1 rounded-md ${badge.bg} ${badge.text}`}
          >
            {w.status}
          </span>
        </div>

       

        <div className="flex gap-2 mt-4 pt-3 border-t border-border">
          <Button
            onClick={() => router.push(`/dashboard/workspaces/${w.id}`)}
            className="rounded-lg"
            size="sm"
          >
            details
          </Button>
        </div>
      </div>
    </div>
  );
}

export function InfoCard({
  title,
  children,
  className,
}: {
  title: string;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "bg-card border border-border rounded-sm p-4",
        className
      )}
    >
      <div className="text-[10px] uppercase tracking-widest text-muted-foreground/50 mb-3">
        {title}
      </div>
      {children}
    </div>
  );
}

export function InfoRow({
  label,
  value,
  mono,
  className,
}: {
  label: string;
  value?: string;
  mono?: boolean;
  className?: string;
}) {
  return (
    <div className="flex justify-between items-center py-1.5 border-b border-border last:border-0">
      <span className="text-[11px] text-muted-foreground">{label}</span>
      <span
        className={cn(
          "text-[11px] text-foreground/80",
          mono && "font-mono text-[10px]",
          className
        )}
      >
        {value ?? "—"}
      </span>
    </div>
  );
}