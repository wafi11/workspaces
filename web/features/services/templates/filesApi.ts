import { api } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

interface Templatefiles {
  id: string;
  template_id: string;
  filename: string;
  sort_order: number;
}

export type EditState = {
  filename: string;
  sort_order: number;
};

export function useUpdateTemplateFiles(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-files", id],
    mutationFn: async (req: EditState) => {
      const response = await api.put(`/templates/files/${id}`, req);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["template-files", templateId],
      });
    },
  });
}
export function useDeleteTemplateFiles(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-files", id],
    mutationFn: async () => {
      const response = await api.delete(`/templates/files/${id}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["template-files", templateId],
      });
    },
  });
}

export function useGetTemplateFiles(id: string) {
  return useQuery({
    queryKey: ["template-files", id],
    queryFn: async () => {
      const response = await api.get<Templatefiles[]>(`/templates/${id}/files`);
      return response.data;
    },
  });
}
