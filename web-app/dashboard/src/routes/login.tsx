import { createFileRoute, Link } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import * as Label from "@radix-ui/react-label";
import { AuthContainer } from "../features/components/AuthContainer";
import { useLogin, useRegisterOrLoginWithGithub } from "@/features/api/auth";
import type { LoginForm } from "@/types";
import { inputStyle } from "@/features/components/InputStyle";
import { ActionBtn } from "@/features/components/ActionButton";

export const Route = createFileRoute("/login")({
  component: RouteComponent,
});

function RouteComponent() {
  const { mutate, isPending } = useLogin();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = (data: LoginForm) => {
    mutate(data);
  };


  return (
    <AuthContainer
      title="Cluster Access"
      subtitle="Otentikasi aman untuk mengelola infrastruktur cloud native Anda."
    >
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-bold text-white tracking-tight">
            Sign In
          </h1>
          <p className="text-xs text-zinc-500">
            Gunakan kredensial akun KubeSpace Anda.
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-4">
            {/* Email Field */}
            <div className="grid gap-2">
              <Label.Root
                className="text-[13px] font-medium text-zinc-300 ml-1"
                htmlFor="email"
              >
                Email Address
              </Label.Root>
              <input
                {...register("email", { required: "Email is required" })}
                id="email"
                type="email"
                placeholder="name@company.com"
                style={inputStyle}
              />
              {errors.email && (
                <span className="text-[10px] text-red-400 ml-1">
                  {errors.email.message}
                </span>
              )}
            </div>

            {/* Password Field */}
            <div className="grid gap-2">
              <div className="flex items-center justify-between px-1">
                <Label.Root
                  className="text-[13px] font-medium text-zinc-300"
                  htmlFor="password"
                >
                  Password
                </Label.Root>
             
              </div>
              <input
                {...register("password", { required: "Password is required" })}
                id="password"
                type="password"
                placeholder="••••••••"
                               style={inputStyle}

              />
              {errors.password && (
                <span className="text-[10px] text-red-400 ml-1">
                  {errors.password.message}
                </span>
              )}
            </div>
          </div>

          <ActionBtn label="Login" variant="default" onClick={() => {}} type="submit" disabled={isPending} className="w-full bg-sidebar-accent/20 py-2"/>
        </form>

        {/* Divider */}
        <div className="relative py-2">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-zinc-800" />
          </div>
          <div className="relative flex justify-center text-[10px] uppercase tracking-widest">
            <span className="bg-[#121214] px-3 text-zinc-600">
              Secure Access
            </span>
          </div>
        </div>

        {/* Social Logins */}
        <div className="grid grid-cols-2 gap-3">
          <button  className="flex items-center justify-center gap-2 rounded-xl border border-zinc-800 bg-transparent h-10 text-xs font-medium text-zinc-300 hover:bg-zinc-900 transition-colors">
            Google
          </button>
          <button onClick={() => useRegisterOrLoginWithGithub()} className="flex items-center justify-center gap-2 rounded-xl border border-zinc-800 bg-transparent h-10 text-xs font-medium text-zinc-300 hover:bg-zinc-900 transition-colors">
            GitHub
          </button>
        </div>

        <p className="text-center text-sm text-zinc-500">
          New to KubeSpace?{" "}
          <Link
            to="/register"
            className="font-semibold text-(--blue-11) hover:text-(--blue-9) hover:underline transition-colors"
          >
            Create Account
          </Link>
        </p>
      </div>
    </AuthContainer>
  );
}
