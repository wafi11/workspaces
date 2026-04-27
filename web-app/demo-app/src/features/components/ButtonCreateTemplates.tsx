import { cn } from "@/lib/utils";

export function ButtonCreateTemplates({ className }: { className?: string }) {
    return (
        <button className={cn(`bg-sidebar-accent/40 px-4 py-1.5 text-sm w-full h-full rounded hover:bg-blue-600 transition duration-300 ${className}`)}>
            <span>+ Create Template</span>
        </button>
    )
}