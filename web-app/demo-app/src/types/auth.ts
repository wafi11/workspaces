import { z } from "zod"


export const loginSchema = z.object({
  email: z
    .string()
    .min(1, "Email wajib diisi")
    .email("Format email tidak valid"),
  password: z
    .string()
    .min(1, "Password wajib diisi")
    .min(8, "Password minimal 8 karakter"),
})


export const registerSchema = z.object({
  name: z
    .string()
    .min(1, "Nama wajib diisi")
    .min(2, "Nama minimal 2 karakter")
    .max(100,"Nama maksimum 100 karakter"),
  username: z
    .string()
    .min(1, "Username wajib diisi")
    .min(3, "Username minimal 3 karakter")
    .max(100,"username maksimum 100 karakter")
    .regex(/^[a-zA-Z0-9_]+$/, "Username hanya boleh huruf, angka, dan underscore"),
  email: z
    .string()
    .min(1, "Email wajib diisi")
    .max(200,"Email maksimum 200 karakter")
    .email("Format email tidak valid"),
  password: z
    .string()
    .min(1, "Password wajib diisi")
    .min(8, "Password minimal 8 karakter")
    .max(15,"Password maksimum 15 karakter")
})

export type Login = z.infer<typeof loginSchema>
export type Register = z.infer<typeof registerSchema>