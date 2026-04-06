import { api } from "@/lib/api";
import { ApiResponse } from "@/types";
import { Workspaces } from "@/types/workspaces";
import { useQuery } from "@tanstack/react-query";

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
