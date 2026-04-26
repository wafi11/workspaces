import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

interface CardActivityProps {
  className?: string;
  title: string;
  children?: ReactNode;
  body?: ReactNode;
  action?: ReactNode;
}

export function CardActivity({ className, title, children, body, action }: CardActivityProps) {
  return (
    <div
      className={cn("flex flex-col rounded-lg overflow-hidden", className)}
      style={{
        background: "var(--color-app-surface)",
        border: "0.5px solid var(--color-app-border)",
      }}
    >
      {/* Header */}
      <div
        className="flex items-center justify-between px-4 py-3"
        style={{ borderBottom: "0.5px solid var(--color-app-border)" }}
      >
        <h2
          className="text-xs font-medium"
          style={{ color: "var(--color-sidebar-text)" }}
        >
          {title}
        </h2>
        {action && (
          <span
            className="text-[11px] cursor-pointer"
            style={{ color: "var(--color-sidebar-text-muted)" }}
          >
            {action}
          </span>
        )}
        {children}
      </div>

      {/* Body */}
      {body && <div>{body}</div>}
    </div>
  );
}