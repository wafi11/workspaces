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
}