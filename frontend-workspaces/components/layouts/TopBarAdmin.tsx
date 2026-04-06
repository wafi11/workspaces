import { ReactNode } from "react";
import { Badge } from "../ui/badge";

interface TopBarAdminProps {
  title: string;
  description: string;
  children?: ReactNode;
}
export function TopBarAdmin({
  description,
  title,
  children,
}: TopBarAdminProps) {
  return (
    <div className="flex items-center justify-between pb-5 border-b border-border">
      <div className="flex flex-col gap-0.5">
        <div className="flex items-center gap-2">
          <Badge
            variant="outline"
            className="bg-primary/10 text-primary border-primary/20 text-[10px] uppercase font-bold tracking-tighter"
          >
            {title}
          </Badge>
        </div>
        <h1 className="text-xl font-medium">{description}</h1>
      </div>
      <div className="flex gap-2">{children}</div>
    </div>
  );
}
