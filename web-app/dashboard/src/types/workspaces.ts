export type Workspaces = {
  id: string;
  user_id: string;
  template_id: string;
  template_name: string;
  name: string;
  icon: string;
  namespace: string;
  status: string;
  env_vars: Record<string, string>;
  created_at: string;
  updated_at: string;
  url: string;
};

export type WorkspaceSessions = {
  created_at: string;
  icon: string;
  env_vars: Record<string, string> | null;
  next_start_at: string;
  expires_at: string;
  id: string;
  name: string;
  started_at: string;
  status: string;
  stopped_at: string | null;
  timezone: string;
  updated_at: string;
  url: string;
};

export type WorkspaceCollaborations = {
  workspace_id: string;
  workspace_name: string;
  workspace_url: string;
  role: "viewer" | "editor";
  status: string;
  invited_at: string;
  template_name: string;
  template_icon: string;
};
export interface WorkspaceRequest {
  name: string;
  description?: string;
  type_time_duration: "minutes" | "hours";
  time_duration: number;
  template_id: string;
  limit_ram_mb: number;
  limit_cpu_cores: number;
  req_ram_mb: number;
  req_cpu_cores: number;
  env_vars?: Record<string, string>;
}

export interface ContainerMetrics {
  name: string;
  cpu_cores: number;
  memory_mb: number;
}

export interface PodMetrics {
  pod_name: string;
  app_name: string;
  containers: ContainerMetrics[];
  total_cpu_cores: number;
  total_memory_mb: number;
}


export interface PodStorageMetrics {
   volumes:PodVolumesType[]
}


export interface PodVolumesType {
  used_mb: number
  capacity_mb: number
  available_mb: number
}
