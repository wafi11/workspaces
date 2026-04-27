import { cn } from "@/lib/utils";
import { useNavigate } from "react-router-dom";

export function ButtonCreateWorkspace({ className }: { className?: string }) {
    const navigate = useNavigate();
    return (
        <button onClick={() => navigate("/workspaces/create")} className={cn(`bg-sidebar-accent/40 px-4 py-1.5 text-sm w-full h-full rounded hover:bg-blue-600 transition duration-300 ${className}`)}>
            <span>+ Create Workspace</span>
        </button>
    )
}