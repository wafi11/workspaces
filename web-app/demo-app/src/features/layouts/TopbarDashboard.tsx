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
      className={cn("flex items-center justify-between  p-2",className)}
      style={{ borderBottom: "1px solid var(--color-sidebar-border)" }}
    >
        {/* Heading */}
        <h1
          className={cn("text-xl font-semibold tracking-tight leading-tight",classNameFont)}
          style={{ color: "var(--color-sidebar-text-active)" }}
        >
          {title}
        </h1>

      

      {/* Actions */}
      {children}
    </div>
  );
}
