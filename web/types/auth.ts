export type User = {
  created_at: string;
  email: string;
  id: string;
  updated_at: string;
  username: string;
};



export type UserQuota = {
  id: string
maxCpuCores: number
maxRamMb: number
maxStorageGb: number
maxWorkspaces: number
userId: string
used_workspaces : number
used_cpu_cores : number
used_ram_mb : number
used_storage_gb : number
}