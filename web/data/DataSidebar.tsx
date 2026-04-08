import {
  Boxes,
  FileCode2,
  LayoutDashboard,
  LucideProps,
  Settings,
} from "lucide-react";
import { ForwardRefExoticComponent, RefAttributes } from "react";

export interface NavItem {
  label: string;
  icon: ForwardRefExoticComponent<
    Omit<LucideProps, "ref"> & RefAttributes<SVGSVGElement>
  >;
  href: string;
  active?: boolean;
}

export const navItems: NavItem[] = [
  { label: "Dashboard", icon: LayoutDashboard, href: "/dashboard" },
  { label: "Workspaces", icon: Boxes, href: "/dashboard/workspaces" },
  { label: "Templates", icon: FileCode2, href: "/dashboard/templates" },
  { label: "Settings", icon: Settings, href: "/dashboard/settings" },
];
export const navItemsUsers: NavItem[] = [
  { label: "Workspaces", icon: LayoutDashboard, href: "/profile" },
  { label: "Templates", icon: FileCode2, href: "/profile/templates" },
  { label: "Observability", icon: Boxes, href: "/profile/obeservability" },
  { label: "Settings", icon: Settings, href: "/profile/settings" },
];
