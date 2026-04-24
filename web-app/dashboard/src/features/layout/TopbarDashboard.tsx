import type { ReactNode } from "react";
import { cn } from "@/lib/utils"

interface TopBarAdminProps {
  title: string;
  children?: ReactNode;
  className? : string
  classNameFont? : string
}

export function TopbarAdmin({ title, children ,className,classNameFont}: TopBarAdminProps) {
  return (
    <div
      className={cn("flex items-center justify-between px-6 py-4",className)}
      style={{ borderBottom: "1px solid var(--color-sidebar-border)" }}
    >
      <div className="flex justify-between w-full items-center gap-1.5">
        {/* Heading */}
        <h1
          className={cn("text-xl font-semibold tracking-tight leading-tight",classNameFont)}
          style={{ color: "var(--color-sidebar-text-active)" }}
        >
          {title}
        </h1>

      </div>

      {/* Actions */}
      {children}
    </div>
  );
}
