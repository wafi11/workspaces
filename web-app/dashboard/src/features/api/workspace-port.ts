import { api } from "@/lib/api";
import { useMutation, useQuery } from "@tanstack/react-query";

export function useCreateWorkspacePort({workspaceId} : {workspaceId : string}){
    return useMutation({
        mutationKey : ["port",workspaceId],
        mutationFn : async ({data} : {data : {port : number}})  => {
            const req = await api.post(`/workspaces/${workspaceId}/port`,data)
            return req.data
        }
    })
}

export function useGetWorkspacesPort({workspaceId} : {workspaceId : string}){
      return useQuery({
        queryKey : ["port",workspaceId],
        queryFn : async ()  => {
            const req = await api.get(`/workspaces/${workspaceId}/port`)
            return req.data
        }
    })
}