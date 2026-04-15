import { api } from "@/lib/api";
import type { TemplateEditVariable, TemplateVariables } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";



export function useUpdateTemplateVariables( templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-variables"],
    mutationFn: async ({id,req}: {id: string,req : TemplateEditVariable}) => {
      const response = await api.put(`/templates/variables/${id}`, req);
      return response.data;
    },
    onSuccess: () => {
      // Handle success, e.g., refetch the updated data
      queryClient.invalidateQueries({
        queryKey: ["template-variables", templateId],
      });
    },
  });
}
export function useCreateTemplateVariables(templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-variables", templateId],
    mutationFn: async (req: TemplateEditVariable) => {
      const response = await api.post(
        `/templates/${templateId}/variables`,
        req
      );
      return response.data;
    },
    onSuccess: () => {
      // Handle success, e.g., refetch the updated data
      queryClient.invalidateQueries({
        queryKey: ["template-variables", templateId],
      });
    },
  });
}
export function useDeleteTemplateVariables(templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-variables", templateId],
    mutationFn: async (id: string) => {
      const response = await api.delete(`/templates/variables/${id}`);
      return response.data;
    },
    onSuccess: () => {
      // Handle success, e.g., refetch the updated data
      queryClient.invalidateQueries({
        queryKey: ["template-variables", templateId],
      });
    },
  });
}

export function useGetTemplateVariables(id?: string) {
  return useQuery({
    queryKey: ["template-variables", id],
    queryFn: async () => {
      const response = await api.get<TemplateVariables[]>(
        `/templates/${id}/variables`
      );
      return response.data;
    },
    enabled : !!id
  });
}
