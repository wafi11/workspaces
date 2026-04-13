import z from "zod";

export const registerSchema = z.object({
  username: z
    .string()
    .min(4, "username is required.")
    .max(100, "username max 100 length characters")
    .trim(),
  email: z
    .string()
    .email("Enter a valid email address.")
    .max(100, "Email max 100 length characters"),
  password: z.string().min(8, "Password must be at least 8 characters."),
});
export const loginSchema = z.object({
  email: z
    .string()
    .email("Enter a valid email address.")
    .max(100, "Email max 100 length characters"),
  password: z.string().min(8, "Password must be at least 8 characters."),
});

export type RegisterForm = z.infer<typeof registerSchema>;
export type LoginForm = z.infer<typeof loginSchema>;
