import { ButtonCreateTemplates } from "@/features/components/ButtonCreateTemplates";
import { ButtonNotification } from "@/features/components/ButtonNotifications";
import { MainContainer } from "@/features/layouts/MainContainer";
import { TopbarAdmin } from "@/features/layouts/TopbarDashboard";

export function TemplatesPage(){
    return (
        <>
            <MainContainer>
                <TopbarAdmin title="Templates">
                    <div className="flex items-center gap-4">
                        <ButtonNotification />
                        <ButtonCreateTemplates />
                    </div>
                </TopbarAdmin>
            </MainContainer>
        </>
    )
}