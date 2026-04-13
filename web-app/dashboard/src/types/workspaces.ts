export type Workspaces = {
  id: string;
  user_id: string;
  template_id: string;
  name: string;
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

export type WorkspaceRequest = {
  name: string;
  description: string;
  env_vars: Record<string, string>;
};
