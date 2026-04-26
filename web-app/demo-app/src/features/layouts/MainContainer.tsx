import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

interface MainContainerProps {
  className?: string;
  children: ReactNode; 
}

export function MainContainer({ className, children }: MainContainerProps) {
  return (
    <main className={cn("flex flex-col flex-1 min-w-0 h-screen overflow-y-auto", className)}>
      {children}
    </main>
  );
}