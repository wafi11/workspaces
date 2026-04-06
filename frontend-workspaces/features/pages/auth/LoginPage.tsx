"use client";
import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { LoginForm, loginSchema, useLogin } from "@/features/services/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { AuthContainer } from "./AuthContainer";
import { Icons } from "@/components/icons";

export function LoginPage() {
  const { mutate, isPending } = useLogin();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = (values: LoginForm) => {
    mutate(values);
  };

  return (
    <AuthContainer
      title="Access Your Workspace"
      subtitle="Kelola environment development hingga production dengan mudah."
    >
      <div className="space-y-4">
        <div className="space-y-2 text-center md:text-left">
          <h1 className="text-2xl font-bold tracking-tight">Sign In</h1>
          <p className="text-sm text-muted-foreground">
            Selamat datang kembali di platform kami.
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)}>
          <FieldGroup>
            {/* Email Field */}
            <Field>
              <FieldLabel htmlFor="email">Email Address</FieldLabel>
              <Input
                id="email"
                type="email"
                placeholder="wafi@example.com"
                {...register("email")}
                className={`rounded-full ${errors.email ? "border-destructive" : ""} `}
              />
              {/* Menampilkan error secara manual menggunakan gaya Field */}
              {errors.email && (
                <p className="text-[10px] font-medium text-destructive">
                  {errors.email.message}
                </p>
              )}
            </Field>

            {/* Password Field */}
            <Field>
              <div className="flex items-center justify-between">
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <Link
                  href="#"
                  className="text-[11px] text-primary hover:underline"
                >
                  Lupa password?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                {...register("password")}
                className={`rounded-full ${errors.password ? "border-destructive" : ""} `}
              />
              {errors.password && (
                <p className="text-[10px] font-medium text-destructive">
                  {errors.password.message}
                </p>
              )}
            </Field>

            <Button
              type="submit"
              className="w-full rounded-full font-bold"
              disabled={isPending}
            >
              {isPending ? "Connecting..." : "Log In"}
            </Button>
          </FieldGroup>
        </form>

        <div className="relative py-2">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-border" />
          </div>
          <div className="relative flex justify-center text-[10px] uppercase">
            <span className="bg-[#1c1c1c] px-2 text-muted-foreground">Or</span>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Button
            variant="outline"
            type="button"
            className="text-xs rounded-full"
          >
            <Icons.google className="mr-2 h-4 w-4" /> Google
          </Button>
          <Button
            variant="outline"
            type="button"
            className="text-xs rounded-full"
          >
            <Icons.gitHub className="mr-2 h-4 w-4" /> Github
          </Button>
        </div>

        <p className="text-center text-sm text-muted-foreground">
          Belum punya akun?{" "}
          <Link href="/register" className="text-primary hover:underline">
            Daftar
          </Link>
        </p>
      </div>
    </AuthContainer>
  );
}
