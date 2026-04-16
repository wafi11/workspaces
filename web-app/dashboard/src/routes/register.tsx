import { useRegister } from "@/features/api/auth";
import type { RegisterForm } from "@/types";
import * as Label from "@radix-ui/react-label";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { AuthContainer } from "../features/components/AuthContainer";
import { ActionBtn } from "@/features/components/ActionButton";
import { inputStyle } from "@/features/components/InputStyle";

export const Route = createFileRoute("/register")({
  component: RouteComponent,
});

function RouteComponent() {
  const { mutate, isPending } = useRegister();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const onSubmit = (data: RegisterForm) => {
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
            Create Account
          </h1>
          <p className="text-xs text-zinc-500">
            Daftarkan akun KubeSpace baru Anda.
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-4">
            {/* Username Field */}
            <div className="grid gap-2">
              <Label.Root
                className="text-[13px] font-medium text-zinc-300 ml-1"
                htmlFor="username"
              >
                Username
              </Label.Root>
              <input
                {...register("username", {
                  required: "Username is required",
                  minLength: { value: 3, message: "Minimal 3 karakter" },
                })}
                id="username"
                type="text"
                placeholder="your-username"
                style={inputStyle}
              />
              {errors.username && (
                <span className="text-[10px] text-red-400 ml-1">
                  {errors.username.message}
                </span>
              )}
            </div>

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
              <Label.Root
                className="text-[13px] font-medium text-zinc-300 ml-1"
                htmlFor="password"
              >
                Password
              </Label.Root>
              <input
                {...register("password", {
                  required: "Password is required",
                  minLength: { value: 8, message: "Minimal 8 karakter" },
                })}
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

         <ActionBtn label="Register" variant="default" onClick={() => {}} type="submit" disabled={isPending} className="w-full bg-sidebar-accent/20 py-2"/>
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
          <button className="flex items-center justify-center gap-2 rounded-xl border border-zinc-800 bg-transparent h-10 text-xs font-medium text-zinc-300 hover:bg-zinc-900 transition-colors">
            Google
          </button>
          <button className="flex items-center justify-center gap-2 rounded-xl border border-zinc-800 bg-transparent h-10 text-xs font-medium text-zinc-300 hover:bg-zinc-900 transition-colors">
            GitHub
          </button>
        </div>

        <p className="text-center text-sm text-zinc-500">
          Sudah punya akun?{" "}
          <Link
            to="/login"
            className="font-semibold text-(--blue-11) hover:text-(--blue-9) hover:underline transition-colors"
          >
            Sign In
          </Link>
        </p>
      </div>
    </AuthContainer>
  );
}