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
