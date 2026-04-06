import { api } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import id from "zod/v4/locales/id.js";

interface TemplateAddOn {
  id: string;
  template_id: string;
  name: string;
  image: string;
  description: string;
}

export type EditState = {
  name: string;
  image: string;
  description: string;
};

export function useCreateTemplateAddOn(templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-add-ons"],
    mutationFn: async (req: EditState) => {
      const response = await api.post(`/templates/${templateId}/add-ons`, req);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["template-add-ons", templateId],
      });
    },
  });
}

export function useUpdateTemplateAddOn(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-add-ons", id],
    mutationFn: async (req: EditState) => {
      const response = await api.put(`/templates/add-ons/${id}`, req);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["template-add-ons", templateId],
      });
    },
  });
}
export function useDeleteTemplateAddOn(id: string, templateId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: ["template-add-ons", id],
    mutationFn: async () => {
      const response = await api.delete(`/templates/add-ons/${id}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["template-add-ons", templateId],
      });
    },
  });
}

export function useGetTemplateAddOns(id: string) {
  return useQuery({
    queryKey: ["template-add-ons", id],
    queryFn: async () => {
      const response = await api.get<TemplateAddOn[]>(
        `/templates/${id}/add-ons`
      );
      return response.data;
    },
  });
}
