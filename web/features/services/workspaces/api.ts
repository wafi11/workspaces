import { api } from "@/lib/api";
import { ApiResponse } from "@/types";
import { Workspaces } from "@/types/workspaces";
import { useMutation, useQuery } from "@tanstack/react-query";

export function useGetListWorkspace() {
  return useQuery({
    queryKey: ["workspaces"],
    queryFn: async () => {
      const request = await api.get<ApiResponse<Workspaces[]>>("/workspaces");
      return request.data;
    },
  });
}

export function useGetWorkspace(id: string) {
  return useQuery({
    queryKey: ["workspaces", id],
    queryFn: async () => {
      const request = await api.get<ApiResponse<Workspaces>>(
        `/workspaces/${id}`
      );
      return request.data;
    },
  });
}
export function useGetListWorkspaceForm() {
  return useQuery({
    queryKey: ["workspaces-form"],
    queryFn: async () => {
      const request = await api.get<{id : string,name : string}[]>(
        `/workspaces/form`
      );
      return request.data;
    },
  });
}

export function useCreateAddonWorkspaces(){
  return useMutation({
    mutationKey : ["workspace-addon"],
    mutationFn : async(req : {workspace_id : string,template_addon_id: string,config : {
      key : string,
      value : string
    }[]})  => {
      const request = await api.post(`/workspaces/${req.workspace_id}/add-ons`,req)
      return request.data
    }
  })
}

export function useGetWorkspacesUsers(){
  return useQuery({
    queryKey : ["workspace-users"],
    queryFn : async ()  => {
      const request = await api.get("/workspaces/user")
      return request.data
    }
  })
}