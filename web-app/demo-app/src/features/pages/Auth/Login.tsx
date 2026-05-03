import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AuthContainer } from "@/features/components/AuthContainer";
import { useLoginForm } from "@/features/hooks/useAuth";
import { Eye, EyeOff } from "lucide-react";
import { useState } from "react";

export function LoginPage() {
  const [isShowPassword, setIsShowPassword] = useState<boolean>(false);
  const { errors, handleSubmit, isPending, register } = useLoginForm();

  return (
    <AuthContainer>
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Email Field */}
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

        {/* Password Field */}
        <div className="flex flex-col gap-1.5">
          {/* Label + Forgot Password sejajar */}
          <div className="flex items-center justify-between">
            <label className="text-sm font-medium text-foreground">
              Password
            </label>
            <a
              href="/forgot-password"
              className="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              Forgot password?
            </a>
          </div>
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

        {/* Submit */}
        <Button type="submit" disabled={isPending} className="w-full">
          {isPending ? "Signing in..." : "Sign in"}
        </Button>

        {/* Register link — di bawah button */}
        <p className="text-center text-sm text-muted-foreground">
          Don't have an account?{" "}
          <a
            href="/register"
            className="text-foreground font-medium hover:underline transition-all"
          >
            Register
          </a>
        </p>
      </form>
    </AuthContainer>
  );
}