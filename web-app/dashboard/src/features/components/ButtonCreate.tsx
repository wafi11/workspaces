import { cn } from "@/lib/utils";
import { useRouter } from "@tanstack/react-router";

interface ButtonCreateProps {
  label : string
  className? : string
  to : string
}

export function ButtonCreate({label,to,className} : ButtonCreateProps) {
  const {navigate} = useRouter()
  return (
      <button
          onClick={() => navigate({ href: `/${to}/create` })}
          className={cn("px-3 py-1.5 w-fit  rounded-md text-sm font-medium transition-colors cursor-pointer",className)}
          style={{
            background: "var(--color-sidebar-surface)",
            border: "1px solid var(--color-sidebar-border)",
            color: "var(--color-sidebar-text-active)",
          }}
        >
          {label}
        </button>
  );
}
