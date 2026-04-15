import { api } from "@/lib/api";
import type { WorkspaceRequest, Workspaces, WorkspaceSessions } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const workspaceUserKeys = {
  all: ["workspace-users"] as const,
};

export function useGetListWorkspace() {
  return useQuery({
    queryKey: ["workspaces"],
    queryFn: async () => {
      const request = await api.get<Workspaces[]>("/workspaces");
      return request.data;
    },
  });
}
export function useCreateWorkspaces() {
  return useMutation({
    mutationKey: ["workspaces"],
    mutationFn: async (req: WorkspaceRequest) => {
      const request = await api.post<{workspace : Workspaces}>("/workspaces", req);
      return request.data;
    },
  });
}

export function useGetWorkspace(id: string) {
  return useQuery({
    queryKey: ["workspaces", id],
    queryFn: async () => {
      const request = await api.get<Workspaces>(`/workspaces/${id}`);
      return request.data;
    },
  });
}
export function useGetListWorkspaceForm() {
  return useQuery({
    queryKey: ["workspaces-form"],
    queryFn: async () => {
      const request =
        await api.get<{ id: string; name: string }[]>(`/workspaces/form`);
      return request.data;
    },
  });
}

export function useCreateAddonWorkspaces() {
  return useMutation({
    mutationKey: ["workspace-addon"],
    mutationFn: async (req: {
      workspace_id: string;
      template_addon_id: string;
      config: {
        key: string;
        value: string;
      }[];
    }) => {
      const request = await api.post(
        `/workspaces/${req.workspace_id}/add-ons`,
        req,
      );
      return request.data;
    },
  });
}

export function useGetAddonWorkspaces(workspaceId: string) {
  return useQuery({
    queryKey: ["workspace-addon", workspaceId],
    queryFn: async () => {
      const request = await api.get(`/workspaces/${workspaceId}/add-ons`);
      return request.data;
    },
  });
}

export function useGetWorkspacesUsers() {
  return useQuery({
    queryKey: workspaceUserKeys.all,
    queryFn: async () => {
      const request = await api.get<WorkspaceSessions[]>("/workspaces/user");
      return request.data;
    },
    staleTime: 0,
  });
}

export function useUpdateStatusWorkspace(id: string, status: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["workspace", id],
    mutationFn: async () => {
      const req = await api.patch(`/workspaces/${id}/${status}`);
      return req.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["workspace", id] });
    },
  });
}
