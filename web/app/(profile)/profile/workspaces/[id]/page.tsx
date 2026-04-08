import { WorkspaceDetailsUser } from "@/features/pages/profile/workspaces/WorkspaceDetails"

export default async function Page({params} : {params : {id : string}}){
    const {id}  = await params
    return (
        <WorkspaceDetailsUser id={id}/>
    )
}