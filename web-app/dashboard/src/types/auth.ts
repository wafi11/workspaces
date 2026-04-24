import type { Role } from "./nav-items";

export type User = {
  created_at: string;
  email: string;
  role : Role
  terminal_url: string;
  id: string;
  avatar_url: string
  updated_at: string;
  username: string;
};

export type UserQuota = {
  id: string;
  maxCpuCores: number;
  maxRamMb: number;
  maxStorageGb: number;
  maxWorkspaces: number;
  userId: string;
  used_workspaces: number;
  used_cpu_cores: number;
  used_ram_mb: number;
  used_storage_gb: number;
};



export interface PAT {
  id: string
  name: string
  expires_at: string
  last_used_at: string | null
  created_at: string
}

export type ProviderUser = {
  name : string
  provider_id : string
}