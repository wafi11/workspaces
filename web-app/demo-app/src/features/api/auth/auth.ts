import { api } from "@/lib/api";
import type { Login, Register } from "@/types";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

export function useLoginApi() {
  const navigate = useNavigate();
  return useMutation({
    mutationKey: ["login"],
    mutationFn: async (data: Login) => {
      const req = await api.post("/auth/login", data);
      return req.data;
    },
    onSuccess: () => {
      navigate("/home");
    },
  });
}

export function useRegisterApi() {
  const navigate = useNavigate();
  return useMutation({
    mutationKey: ["register"],
    mutationFn: async (data: Register) => {
      const req = await api.post("/auth/register", data);
      return req.data;
    },
    onSuccess: () => {
      navigate("/home");
    },
  });
}