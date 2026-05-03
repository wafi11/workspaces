import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AuthContainer } from "@/features/components/AuthContainer";
import { useRegisterForm } from "@/features/hooks/useAuth";
import { Eye, EyeOff } from "lucide-react";
import { useState } from "react";

export function RegisterPage() {
  const [isShowPassword, setIsShowPassword] = useState<boolean>(false);
  const { errors, submit, isPending, register } = useRegisterForm();

  return (
    <AuthContainer>
      <form onSubmit={submit} className="space-y-6">
        
    
        <div className="flex flex-col gap-1.5">
          <label className="text-sm font-medium text-foreground">Name</label>
          <Input
            {...register("name")}
            type="text"
            placeholder="joko samudra"
            className="border border-blue-300/20 bg-muted/30 focus:border-blue-400/60"
          />
          {errors.name && (
            <p className="text-xs text-destructive">{errors.name.message}</p>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <label className="text-sm font-medium text-foreground">Username</label>
          <Input
            {...register("username")}
            type="text"
            placeholder="jokosamudra11"
            className="border border-blue-300/20 bg-muted/30 focus:border-blue-400/60"
          />
          {errors.username && (
            <p className="text-xs text-destructive">{errors.username.message}</p>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <label className="text-sm font-medium text-foreground">Email</label>
          <Input
            {...register("email")}
            type="email"
            placeholder="you@example.com"
            className="border border-blue-300/20 bg-muted/30 focus:border-blue-400/60"
          />
          {errors.email && (
            <p className="text-xs text-destructive">{errors.email.message}</p>
          )}
        </div>

        
        <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-foreground">
              Password
            </label>
          <div className="relative">
            <Input
              {...register("password")}
              type={isShowPassword ? "text" : "password"}
              placeholder="••••••••••••"
              className="border border-blue-300/20 bg-muted/30 pr-10 focus:border-blue-400/60"
            />
            <button
              type="button"
              onClick={() => setIsShowPassword((prev) => !prev)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
              tabIndex={-1}
            >
              {isShowPassword ? (
                <EyeOff className="size-4" />
              ) : (
                <Eye className="size-4" />
              )}
            </button>
          </div>
          {errors.password && (
            <p className="text-xs text-destructive">
              {errors.password.message}
            </p>
          )}
        </div>

        
        <Button type="submit" disabled={isPending} className="w-full">
          {isPending ? "Signing up..." : "Sign up"}
        </Button>

        <p className="text-center text-sm text-muted-foreground">
          Do have an account?{" "}
          <a
            href="/login"
            className="text-foreground font-medium hover:underline transition-all"
          >
            Login
          </a>
        </p>
      </form>
    </AuthContainer>
  );
}