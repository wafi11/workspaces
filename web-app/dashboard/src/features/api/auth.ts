"use client";
import { storage } from "@/hooks";
import { api } from "@/lib/api";
import type { User, UserQuota } from "@/types";
import type { LoginForm, RegisterForm } from "@/types/validation-auth";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useRouter } from "@tanstack/react-router";

export function useRegister() {
  const router = useRouter();
  return useMutation({
    mutationKey: ["register"],
    mutationFn: async (data: RegisterForm) => {
      const req = await api.post<null>("/auth/register", data);
      return req.data;
    },
    onSuccess: () => {
      //   router.(LOGIN_URL);
    },
  });
}

export function useLogin() {
  const router = useRouter();

  return useMutation({
    mutationKey: ["login"],
    mutationFn: async (data: LoginForm) => {
      const req = await api.post<{
        access_token: string;
        refresh_token: string;
      }>("/auth/login", data);
      const { access_token, refresh_token } = req.data;
      storage.setAccessToken(access_token);
      storage.setRefreshToken(refresh_token);

      return req.data;
    },
    onSuccess: () => {
      //   router.push("/dashboard");
    },
  });
}

export function useLogout() {
  const router = useRouter();

  return useMutation({
    mutationKey: ["logout"],
    mutationFn: async () => {
      await api.post("/auth/logout", {
        refresh_token: storage.getRefreshToken(),
      });
    },
    onSuccess: () => {
      storage.clear();
      //   router.push("/login");
    },
    onError: () => {
      // Tetap logout di client meski server error
      storage.clear();
      //   router.push("/login");
    },
  });
}

export function useProfile() {
  return useQuery({
    queryKey: ["profile"],
    queryFn: async () => {
      const req = await api.get<User>("/users");
      return req.data;
    },
  });
}

export function useProfileQuota() {
  return useQuery({
    queryKey: ["profile-quota"],
    queryFn: async () => {
      const req = await api.get<UserQuota>("/users/quota");
      return req.data;
    },
    staleTime: 0,
  });
}
