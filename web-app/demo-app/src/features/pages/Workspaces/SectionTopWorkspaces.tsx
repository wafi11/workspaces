import { ButtonCreateWorkspace } from "@/features/components/ButtonCreateWorskspace";
import {  WorkspaceStatusFilter } from "@/features/components/FilterStatus";
import { SearchBar } from "@/features/components/SearchBar";
import { useSearchWorkspaces } from "@/features/hooks/SearchHooks";

export function SectionTopWorkspaces(){
    const {filteredWorkspaces,searchTerm,setSearchTerm} = useSearchWorkspaces({ workspacesData: [] })
    return (
        <div className="flex items-center justify-between gap-2">
            <SearchBar onChange={setSearchTerm} value={searchTerm} placeholder="Search Workspaces" />  
            <WorkspaceStatusFilter value="running" onChange={() => {}} />  
            <ButtonCreateWorkspace className="px-3 py-1.5 hidden sm:block bg-blue-800 hover:bg-blue-600" />
        </div>
    )
}
