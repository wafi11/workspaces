"use client";
import { cn } from "@/lib/utils";
import { Check } from "lucide-react";


export function StepBar({
  steps,
  current,
  onGoTo,
}: {
  steps: readonly string[];
  current: number;
  onGoTo: (i: number) => void;
}) {
  return (
    <div className="flex items-center border-b border-border my-4 bg-muted/30">
      {steps.map((label, i) => {
        const done = i < current;
        const active = i === current;
        return (
          <button
            key={label}
            onClick={() => done && onGoTo(i)}
            disabled={!done}
            className={cn(
              "flex-1 flex flex-col items-center gap-1 py-3 text-xs transition-colors",
              active && "text-primary font-semibold",
              done &&
                "text-muted-foreground hover:text-foreground cursor-pointer",
              !done && !active && "text-muted-foreground/40 cursor-default"
            )}
          >
            <span
              className={cn(
                "flex h-5 w-5 items-center justify-center rounded-full text-[10px] font-bold",
                active && "bg-primary text-primary-foreground",
                done && "bg-muted-foreground/30 text-foreground",
                !done && !active && "bg-muted text-muted-foreground/40"
              )}
            >
              {done ? <Check className="h-3 w-3" /> : i + 1}
            </span>
            <span className="hidden sm:block">{label}</span>
          </button>
        );
      })}
    </div>
  );
}