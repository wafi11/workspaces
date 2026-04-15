import { api } from "@/lib/api";
import type { Templatefiles, TemplateFilesEdit } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";


export function useUpdateTemplateFiles(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-files", id],
    mutationFn: async (req: TemplateFilesEdit) => {
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
export function useCreateTemplateFiles(templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-files", templateId],
    mutationFn: async (req: TemplateFilesEdit) => {
      const response = await api.post(`/templates/${templateId}/files`, req);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["template-files", templateId],
      });
    },
  });
}
export function useDeleteTemplateFiles(templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-files", templateId],
    mutationFn: async (id: string) => {
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
