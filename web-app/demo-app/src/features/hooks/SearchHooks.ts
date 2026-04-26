import type { Workspace } from "@/types";
import { useEffect, useState } from "react";




export function useSearchWorkspaces({ workspacesData }: { workspacesData: Workspace[] }) {
    const [searchTerm, setSearchTerm] = useState("");
    const [filterStatus, setFilterStatus] = useState("all");
    const [filteredWorkspaces, setFilteredWorkspaces] = useState(workspacesData);

    useEffect(() => {
        if (searchTerm === "") {
            setFilteredWorkspaces(workspacesData);
        } else {
            const filtered = workspacesData.filter((workspace) =>
                workspace.name.toLowerCase().includes(searchTerm.toLowerCase())
            );
            setFilteredWorkspaces(filtered);
        }
    }, [searchTerm]);

    return { searchTerm, setSearchTerm, filteredWorkspaces };
}