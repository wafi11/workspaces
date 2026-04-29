import { MainContainer } from "@/features/layouts/MainContainer";
import { FormWorkspaces } from "./FormWorkspace";

export function CreateWorkspacePage(){
    return (
        <MainContainer className="grid grid-cols-2 gap-4">
            <FormWorkspaces />
        </MainContainer>
    )
}