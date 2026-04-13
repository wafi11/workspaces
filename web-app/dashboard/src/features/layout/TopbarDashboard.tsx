import type { ReactNode } from "react";

interface TopBarAdminProps {
  title: string;
  children?: ReactNode;
}

export function TopbarAdmin({ title, children }: TopBarAdminProps) {
  return (
    <div
      className="flex items-center justify-between px-6 py-4"
      style={{ borderBottom: "1px solid var(--color-sidebar-border)" }}
    >
      <div className="flex flex-col gap-1.5">
        {/* Label pill */}

        {/* Heading */}
        <h1
          className="text-xl font-semibold tracking-tight leading-tight"
          style={{ color: "var(--color-sidebar-text-active)" }}
        >
          {title}
        </h1>
      </div>

      {/* Actions */}
      {children && <div className="flex items-center gap-3">{children}</div>}
    </div>
  );
}
