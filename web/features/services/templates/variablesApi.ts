import { api } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

interface TemplateVariables {
  id: string;
  template_id: string;
  key: string;
  default_value: string;
  required: boolean;
  description: string;
}

export type EditState = {
  key: string;
  default_value: string;
  description: string;
};

export function useUpdateTemplateVariables(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-variables", id],
    mutationFn: async (req: EditState) => {
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
export function useDeleteTemplateVariables(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-variables", id],
    mutationFn: async () => {
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

export function useGetTemplateVariables(id: string) {
  return useQuery({
    queryKey: ["template-variables", id],
    queryFn: async () => {
      const response = await api.get<TemplateVariables[]>(
        `/templates/${id}/variables`
      );
      return response.data;
    },
  });
}
