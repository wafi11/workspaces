import type { Register, Login } from "@/types";
import { useForm } from "react-hook-form";
import { useLoginApi, useRegisterApi } from "../api/auth/auth";

export function useRegisterForm() {
  const {mutate,isPending}  = useRegisterApi()
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Register>();

  const onSubmit = (data: Register) => {
    // handle register
    console.log(data);
    mutate(data)
  };

  return { register,isPending, submit :  handleSubmit(onSubmit), errors };
}

export function useLoginForm() {
  const {mutate,isPending}  = useLoginApi()
  
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Login>({
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = (data: Login) => {
    // handle login
    mutate(data)
    console.log(data);
  };

  return { register,isPending, handleSubmit: handleSubmit(onSubmit), errors };
}