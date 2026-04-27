import { MainContainer } from "@/features/layouts/MainContainer";
import { TopbarAdmin } from "@/features/layouts/TopbarDashboard";
import { SectionTopWorkspaces } from "./SectionTopWorkspaces";
import { SectionListWorkspaces } from "./SectionListWorkspaces";
import type { Workspace } from "@/types/workspace";
import { mockWorkspaces } from "@/features/data/dataMockWorkspace";


export function WorkspacesPage(){
    return (
        <MainContainer>
            <TopbarAdmin title="Workspaces">
                <SectionTopWorkspaces />
            </TopbarAdmin>
            {/* main workspaces */}
            <SectionListWorkspaces workspaces={mockWorkspaces as Workspace[]} />
        </MainContainer>
    )
}