import { api } from "@/lib/api";
import { ApiResponse, CreateTemplateRequest, Templates } from "@/types";
import { useMutation, useQuery } from "@tanstack/react-query";

export function useCreateTemplates() {
  return useMutation({
    mutationKey: ["templates"],
    mutationFn: async (req: CreateTemplateRequest) => {
      const request = await api.post("/templates", req);
      return request.data;
    },
  });
}

export function useGetTemplates() {
  return useQuery({
    queryKey: ["list-templates"],
    queryFn: async () => {
      const req = await api.get<ApiResponse<Templates[]>>("/templates");
      return req.data;
    },
  });
}
export function useGetTemplateDetails(templateId: string) {
  return useQuery({
    queryKey: ["template", templateId],
    queryFn: async () => {
      const req = await api.get<ApiResponse<Templates>>(
        `/templates/${templateId}`
      );
      return req.data;
    },
  });
}

export function useUpdateTemplates(id: string) {
  return useMutation({
    mutationKey: ["templates", id],
    mutationFn: async (data: Partial<Templates>) => {
      const req = await api.put(`/templates/${id}`, data);
      return req.data;
    },
  });
}

export function DeleteTemplates(id: string) {
  return useMutation({
    mutationKey: ["templates", id],
    mutationFn: async () => {
      const req = await api.delete(`/templates/${id}`);
      return req.data;
    },
  });
}
