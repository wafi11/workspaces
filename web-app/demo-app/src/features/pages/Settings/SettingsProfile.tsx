import { ProfileSection } from "@/components/profile/SectionProfile"
import { useGetProfile } from "@/features/api/user/profile"
import { MainContainer } from "@/features/layouts/MainContainer"

export function SettingsProfile(){
    const {data}  = useGetProfile()
    console.log(data)
    if (!data) {
        return <></>
    }
    return (
        <MainContainer>
            <ProfileSection user={data}/>
        </MainContainer>
    )
}
