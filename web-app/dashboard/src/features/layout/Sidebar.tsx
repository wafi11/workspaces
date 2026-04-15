import { useState } from "react";
import { Link, useRouterState } from "@tanstack/react-router";
import { Avatar, Tooltip } from "@radix-ui/themes";
import type { Role } from "../../types";
import { NAV } from "../../data";
import { Ico } from "../components/icons";

interface SidebarProps {
  role?: Role;
  userName?: string;
  userEmail?: string;
  onLogout?: () => void;
}

export function Sidebar({
  role = "user",
  userName = "Wafi",
  userEmail = "wafi@wfdnstore.online",
  onLogout,
}: SidebarProps) {
  const [collapsed, setCollapsed] = useState(false);
  const { location } = useRouterState();
  const path = location.pathname;

  const items = NAV.filter((i) => i.roles.includes(role));

  const hoverOn = (e: React.MouseEvent<HTMLElement>, active = false) => {
    if (!active) {
      e.currentTarget.style.background = "var(--color-sidebar-surface)";
      e.currentTarget.style.color = "var(--color-sidebar-text-active)";
    }
  };
  const hoverOff = (e: React.MouseEvent<HTMLElement>, active = false) => {
    if (!active) {
      e.currentTarget.style.background = "transparent";
      e.currentTarget.style.color = "var(--color-sidebar-text)";
    }
  };

  return (
    <aside
      className={[
        "flex flex-col relative min-h-screen h-full shrink-0",
        "transition-[width] duration-200 ease-in-out",
        collapsed ? "w-16" : "w-56",
      ].join(" ")}
      style={{
        background: "var(--color-sidebar-bg)",
        borderRight: "1px solid var(--color-sidebar-border)",
      }}
    >
       
      {/* Logo Area */}
<div
  className={[
    "flex items-center h-14 shrink-0",
    collapsed ? "justify-center" : "px-4 gap-3",
  ].join(" ")}
  style={{ borderBottom: "1px solid var(--color-sidebar-border)" }}
>
  {/* TOMBOL COLLAPSE - Posisi nempel di garis */}
  <button
    onClick={() => setCollapsed((c) => !c)}
    className={[
      "absolute z-50 flex items-center justify-center transition-all",
      "top-4 -right-3 w-6 h-6 rounded-full border shadow-sm cursor-pointer",
    ].join(" ")}
    style={{ 
      color: "var(--color-sidebar-text-active)",
      background: "var(--color-sidebar-bg)", // Harus sama dengan background sidebar agar seamless
      borderColor: "var(--color-sidebar-border)"
    }}
  >
    <Ico
      d2="12"
      d={
        collapsed
          ? "M9 5l7 7-7 7" // Icon panah ke kanan saja (lebih clean)
          : "M15 19l-7-7 7-7" // Icon panah ke kiri
      }
    />
  </button>

  {/* Logo Icon & Text */}
  <div
    className="w-6 h-6 rounded-md flex items-center justify-center shrink-0"
    style={{ background: "var(--color-sidebar-text-active)" }}
  >
    {/* SVG Logo Anda */}
  </div>
  {!collapsed && (
    <span className="text-sm font-medium tracking-tight truncate" style={{ color: "var(--color-sidebar-text-active)" }}>
      Workspaces
    </span>
  )}
</div>

      {/* Nav */}
      <nav className="flex-1 py-2 flex flex-col gap-0.5 overflow-hidden">
        {items.map((item) => {
          const active = path.startsWith(item.to);

          const linkEl = (
            <Link
              to={item.to}
              className={[
                "flex items-center gap-2.5 mx-2 px-3 py-2 rounded-md text-sm",
                "transition-colors duration-100 relative",
                collapsed ? "justify-center" : "",
              ].join(" ")}
              style={{
                background: active
                  ? "var(--color-sidebar-surface)"
                  : "transparent",
                color: active
                  ? "var(--color-sidebar-text-active)"
                  : "var(--color-sidebar-text)",
              }}
              onMouseEnter={(e) => hoverOn(e, active)}
              onMouseLeave={(e) => hoverOff(e, active)}
            >
              {active && (
                <span
                  className="absolute left-0 top-2 bottom-2 w-0.5 rounded-r-full"
                  style={{ background: "var(--color-sidebar-accent)" }}
                />
              )}
              {item.icon}
              {!collapsed && (
                <>
                  <span className="flex-1 truncate">{item.label}</span>
                  {item.badge && (
                    <span
                      className="text-[10px] px-1.5 py-px rounded leading-none"
                      style={{
                        color: "var(--color-sidebar-text)",
                        background: "var(--color-sidebar-border)",
                        border: "1px solid var(--color-app-border)",
                      }}
                    >
                      {item.badge}
                    </span>
                  )}
                </>
              )}
            </Link>
          );

          return collapsed ? (
            <Tooltip key={item.to} content={item.label} side="right">
              {linkEl}
            </Tooltip>
          ) : (
            <div key={item.to}>{linkEl}</div>
          );
        })}
      </nav>

      {/* Divider */}
      <div
        className="mx-3 shrink-0"
        style={{ height: "1px", background: "var(--color-sidebar-border)" }}
      />

      {/* User */}
      <div
        className={[
          "flex items-center gap-2.5 mx-2 px-3 py-3 rounded-md",
          collapsed ? "justify-center" : "",
        ].join(" ")}
      >
        <Avatar
          size="1"
          fallback={userName.slice(0, 2).toUpperCase()}
          radius="small"
          className="shrink-0"
        />
        {!collapsed && (
          <div className="flex flex-col min-w-0 flex-1">
            <span
              className="text-xs font-medium truncate"
              style={{ color: "var(--color-sidebar-text-active)" }}
            >
              {userName}
            </span>
            <span
              className="text-[11px] truncate leading-tight"
              style={{ color: "var(--color-sidebar-text)" }}
            >
              {userEmail}
            </span>
            <span
              className="text-[10px] mt-1 px-1.5 py-px rounded w-fit leading-none"
              style={
                role === "admin"
                  ? {
                      color: "#a16207",
                      background: "#1c1400",
                      border: "1px solid #292000",
                    }
                  : {
                      color: "var(--color-sidebar-text)",
                      background: "var(--color-sidebar-surface)",
                      border: "1px solid var(--color-sidebar-border)",
                    }
              }
            >
              {role}
            </span>
          </div>
        )}
      </div>

      {/* Logout + Collapse */}
      <div className="px-2 pb-3 flex flex-col gap-0.5">
        {collapsed ? (
          <Tooltip content="Logout" side="right">
            <button
              onClick={onLogout}
              className="flex items-center justify-center p-2 rounded-md w-full transition-colors"
              style={{ color: "var(--color-sidebar-text)" }}
              onMouseEnter={(e) => hoverOn(e)}
              onMouseLeave={(e) => hoverOff(e)}
            >
              <Ico d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" />
            </button>
          </Tooltip>
        ) : (
          <button
            onClick={onLogout}
            className="flex items-center gap-2.5 px-3 py-2 rounded-md text-sm w-full transition-colors"
            style={{ color: "var(--color-sidebar-text)" }}
            onMouseEnter={(e) => hoverOn(e)}
            onMouseLeave={(e) => hoverOff(e)}
          >
            <Ico d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" />
            <span>Logout</span>
          </button>
        )}

      
      </div>
    </aside>
  );
}
