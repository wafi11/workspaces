import { ButtonCreateWorkspace } from '@/features/components/ButtonCreateWorskspace'
import { ButtonNotification } from '@/features/components/ButtonNotifications'
import { MainContainer } from '@/features/layouts/MainContainer'
import { TopbarAdmin } from '@/features/layouts/TopbarDashboard'
import { SectionStatsCard } from './SectionStatsCard'
import { SectionActivity } from './SectionActivity'

export default function HomePage() {
    return (
        <>
            <MainContainer>
                <TopbarAdmin title="Home">
                    <div className="flex items-center gap-4">
                        <ButtonNotification />
                        <ButtonCreateWorkspace />
                    </div>
                </TopbarAdmin>
                <SectionStatsCard />
                <SectionActivity />
            </MainContainer>
        </>
    )
}
