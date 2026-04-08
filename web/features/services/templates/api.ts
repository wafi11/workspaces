import { api } from "@/lib/api";
import { CreateTemplateRequest, TemplateForm, Templates } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

// Centralized query keys
export const templateKeys = {
  all: ["templates"] as const,
  lists: () => [...templateKeys.all, "list"] as const,
  detail: (id: string) => [...templateKeys.all, "detail", id] as const,
  form: (id: string) => [...templateKeys.all, "form", id] as const,
}

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