"use client"
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin"
import { useProfile, useProfileQuota } from "@/features/services/auth"
import { useGetWorkspacesUsers } from "@/features/services/workspaces/api"
import { SectionUserQuota } from "./SectionUserQuota"
import { SectionWorkspaces } from "./SectionWorkspaces" // Import ini
import { getSystemGreeting } from "@/utils"

export function ProfilePage(){
    const { data: profile } = useProfile()
    const { data: dataQuota } = useProfileQuota()
    const { data: workspaceData } = useGetWorkspacesUsers()
    
    const greeting = getSystemGreeting()

    return (
        <div className="container mx-auto pb-10">
            <TopBarAdmin 
                description={`${greeting.sub} Welcome back, ${profile?.username || 'User'}.`} 
                title={greeting.text}
            />
            
            <div className="space-y-8 mt-6">
                {/* Section Quota (Paling Atas sebagai Ringkasan) */}
                <SectionUserQuota data={dataQuota} />
                
                {/* Section Workspaces (Detail Resource) */}
                <SectionWorkspaces data={workspaceData} />
            </div>
        </div>
    )
}