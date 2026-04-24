import { api } from "@/lib/api";
import type { WorkspaceCollaborations } from "@/types";
import { useMutation, useQuery } from "@tanstack/react-query";



export function useAddCollaboration({wsId} :{ wsId : string}){
    return useMutation({
        mutationKey : ["workspace-collaboration",wsId],
        mutationFn : async ({data} : {data : {email : string,role : string}})  => {
            const req = await api.post(`/workspaces/${wsId}/invite`,data)
            return req.data
        },
       
    })
}


export function useGetWorkspaceCollaboration(){
    return useQuery({
        queryKey : ["workspace-collaborations"],
        queryFn : async () => {
            const req = await api.get<WorkspaceCollaborations[]>("/workspaces/collaboration/teams")
            return req.data
        }
    })
}


export function useAcceptOrDenied(){
    return useMutation({
        mutationKey : ["notification-permissions"],
        mutationFn : async ({data} : {data : {notification_id : string,types : string}})  => {
            console.log(data)
            const req = await api.post(`/workspaces/collaborator/permissions`,data)
            return req.data
        }
    })
}