"use client";
import { navItems } from "@/data/DataSidebar";
import { useProfile } from "@/features/services/auth";
import { cn } from "@/lib/utils";
import { LogOut } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

export function SidebarProfile() {
  const pathname = usePathname();

  const { data } = useProfile();
  const user = data;

  const initials = data?.username
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <aside className="fixed left-0 top-0 h-screen w-[260px] border-r bg-sidebar flex flex-col p-4 gap-2">
      {/* Logo */}
      <div className="px-3 py-2 mb-1">
        <p className="font-medium text-[15px]">Dashboard</p>
      </div>

      {/* Nav */}
      <nav className="flex flex-col gap-1 flex-1">
        {navItems.map(({ href, label, icon: Icon }) => (
          <Link
            key={href}
            href={href}
            className={cn(
              "flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors",
              pathname === href
                ? "bg-secondary font-medium text-primary"
                : "text-muted-foreground hover:bg-secondary hover:text-primary"
            )}
          >
            <Icon className="w-4 h-4 flex-shrink-0" />
            {label}
          </Link>
        ))}
      </nav>

      {/* Profile */}
      <div className="border-t pt-3">
        <div className="flex items-center gap-3 px-3 py-2 rounded-lg border">
          <div className="w-9 h-9 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center flex-shrink-0">
            <span className="text-sm font-medium text-blue-700 dark:text-blue-200">
              {initials}
            </span>
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-sm font-medium truncate">{user?.username}</p>
            <p className="text-xs text-muted-foreground truncate">
              {user?.email}
            </p>
          </div>
          <button className="p-1 rounded-md hover:bg-secondary transition-colors">
            <LogOut className="w-4 h-4 text-muted-foreground" />
          </button>
        </div>
      </div>
    </aside>
  );
}
