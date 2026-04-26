import { MainContainer } from "@/features/layouts/MainContainer";
import { TopbarAdmin } from "@/features/layouts/TopbarDashboard";
import { SectionTopWorkspaces } from "./SectionTopWorkspaces";

export function WorkspacesPage(){
    return (
        <MainContainer>
            <TopbarAdmin title="Workspaces">
                <SectionTopWorkspaces />
            </TopbarAdmin>
        </MainContainer>
    )
}