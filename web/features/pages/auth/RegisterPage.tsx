"use client";

import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  RegisterForm,
  registerSchema,
  useRegister,
} from "@/features/services/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { AuthContainer } from "./AuthContainer";
import { Icons } from "@/components/icons";

export function RegisterPage() {
  const { mutate, isPending } = useRegister();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const onSubmit = (values: RegisterForm) => {
    mutate(values);
  };

  return (
    <AuthContainer
      title="Access Your Workspace"
      subtitle="Siapkan infrastruktur Anda dan mulai deploy aplikasi pertama Anda hari ini."
    >
      <div className="space-y-2">
        <div className="space-y-2 text-center md:text-left">
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Create Account
          </h1>
          <p className="text-sm text-muted-foreground">
            Bergabunglah dengan platform workspace masa depan.
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)}>
          <FieldGroup className="">
            {/* Username Field */}
            <Field>
              <FieldLabel htmlFor="username">Username</FieldLabel>
              <Input
                id="username"
                type="text"
                placeholder="wafi_dev"
                {...register("username")}
                className={`rounded-full ${errors.username ? "border-destructive" : ""} `}
              />
              {errors.username && (
                <p className="text-[10px] font-medium text-destructive">
                  {errors.username.message}
                </p>
              )}
            </Field>

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
              {errors.email && (
                <p className="text-[10px] font-medium text-destructive">
                  {errors.email.message}
                </p>
              )}
            </Field>

            {/* Password Field */}
            <Field>
              <FieldLabel htmlFor="password">Password</FieldLabel>
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
              className="w-full rounded-full font-bold mt-2"
              disabled={isPending}
            >
              {isPending ? "Creating Account..." : "Sign Up"}
            </Button>
          </FieldGroup>
        </form>

        <div className="relative py-2">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-border" />
          </div>
          <div className="relative flex justify-center text-[10px] uppercase">
            <span className="bg-[#1c1c1c] px-2 text-muted-foreground">
              Or register with
            </span>
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
          Sudah punya akun?{" "}
          <Link
            href="/login"
            className="text-primary hover:underline font-medium"
          >
            Masuk Sekarang
          </Link>
        </p>
      </div>
    </AuthContainer>
  );
}
