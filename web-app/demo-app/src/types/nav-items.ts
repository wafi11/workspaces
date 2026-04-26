// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------
export type Role = "user" | "admin"

export interface NavItem {
  label: string;
  to: string;
  icon: React.ReactNode;
  roles: Role[];
  badge?: string;
}
