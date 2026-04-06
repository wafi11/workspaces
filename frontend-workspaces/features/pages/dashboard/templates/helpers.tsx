import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export function Field({
  label,
  error,
  children,
}: {
  label: string;
  error?: string;
  children?: React.ReactNode;
}) {
  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-medium text-muted-foreground">
        {label}
      </Label>
      {children}
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  );
}

export function EmptyState({ label }: { label: string }) {
  return (
    <p className="text-xs text-muted-foreground text-center py-6 border border-dashed border-border rounded-lg">
      {label}
    </p>
  );
}

export function ReviewRow({
  label,
  value,
  mono,
}: {
  label: string;
  value: React.ReactNode;
  mono?: boolean;
}) {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-border pb-3 last:border-0 last:pb-0">
      <span className="text-muted-foreground shrink-0 w-24">{label}</span>
      <span className={cn("text-right break-all", mono && "font-mono text-xs")}>
        {value}
      </span>
    </div>
  );
}
