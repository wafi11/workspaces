import { api } from "@/lib/api";
import type {
  CreateTemplateRequest,
  EditState,
  TemplateAddOn,
  TemplateForm,
  Templates,
} from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

// Centralized query keys
export const templateKeys = {
  all: ["templates"] as const,
  lists: () => [...templateKeys.all, "list"] as const,
  detail: (id: string) => [...templateKeys.all, "detail", id] as const,
  form: (id: string) => [...templateKeys.all, "form", id] as const,
  workspace_form: () => [...templateKeys.all, "workspace-form"] as const,
};

export function useTemplates() {
  return useQuery({
    queryKey: templateKeys.lists(),
    queryFn: async () => {
      const { data } = await api.get<Templates[]>("/templates");
      return data;
    },
  });
}

export function useTemplateDetails(id: string) {
  return useQuery({
    queryKey: templateKeys.detail(id),
    queryFn: async () => {
      const { data } = await api.get<Templates>(`/templates/${id}`);
      return data;
    },
    enabled: !!id,
  });
}

export function useTemplateForm(id: string) {
  return useQuery({
    queryKey: templateKeys.form(id),
    queryFn: async () => {
      const { data } = await api.get<TemplateForm>(`/templates/${id}/form`);
      return data;
    },
    enabled: !!id,
  });
}

export function useFindTemplateWorkspaceForm() {
  return useQuery({
    queryKey: templateKeys.workspace_form(),
    queryFn: async () => {
      const { data } = await api.get<
        { id: string; name: string; icon: string }[]
      >(`/templates/workspace/form`);
      return data;
    },
  });
}

export function useCreateTemplate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (req: CreateTemplateRequest) => {
      const { data } = await api.post<Templates>("/templates", req);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: templateKeys.lists() });
    },
  });
}

export function useUpdateTemplate(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: Partial<Templates>) => {
      const { data } = await api.put<Templates>(`/templates/${id}`, body);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: templateKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: templateKeys.lists() });
    },
  });
}

export function useDeleteTemplate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/templates/${id}`);
      return id;
    },
    onSuccess: (id) => {
      queryClient.invalidateQueries({ queryKey: templateKeys.lists() });
      queryClient.removeQueries({ queryKey: templateKeys.detail(id) });
    },
  });
}

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
        `/templates/${id}/add-ons`,
      );
      return response.data;
    },
  });
}
