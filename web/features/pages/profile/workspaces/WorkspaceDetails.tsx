"use client"

import { useGetAddonWorkspaces, useGetWorkspace } from "@/features/services/workspaces/api"
import { WorkspaceCard } from "../../dashboard/workspaces/WorkspaceCard"
import { WorkspaceDetails } from "../../dashboard/workspaces/WorkspaceDetails"
import { WorkspaceAddons } from "./WorkspaceAddon"

interface WorkspaceDetails {
    id : string
}
export function WorkspaceDetailsUser({id} : WorkspaceDetails){
    console.log(id)
   const {data}  = useGetAddonWorkspaces(id)
   console.log(data) 
   return (
        <>
            {data && <WorkspaceAddons addons={data} />}
        </>
    )
}