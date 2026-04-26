import { useState } from "react";
import { SidebarLists } from "./SidebarList";
import { useLocation } from "react-router-dom";
import { NAV } from "../data";
import type { Role } from "@/types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronRight, faChevronLeft, faRightFromBracket } from "@fortawesome/free-solid-svg-icons";
import { Avatar } from "../components/Avatar";

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
  const { pathname } = useLocation();
  const path = pathname;

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
        {/* Tombol Collapse */}
        <button
          onClick={() => setCollapsed((c) => !c)}
          className="absolute z-50 flex items-center justify-center transition-all top-4 -right-3 w-6 h-6 rounded-full border shadow-sm cursor-pointer"
          style={{
            color: "var(--color-sidebar-text-active)",
            background: "var(--color-sidebar-bg)",
            borderColor: "var(--color-sidebar-border)",
          }}
        >
          <FontAwesomeIcon icon={collapsed ? faChevronRight : faChevronLeft} className="text-xs" />
        </button>

        {/* Logo */}
        <div
          className="w-6 h-6 rounded-md flex items-center justify-center shrink-0"
          style={{ background: "var(--color-sidebar-text-active)" }}
        />
        {!collapsed && (
          <span
            className="text-sm font-medium tracking-tight truncate"
            style={{ color: "var(--color-sidebar-text-active)" }}
          >
            {userName}
          </span>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 py-2 flex flex-col gap-0.5 overflow-hidden">
        {items.map((item) => {
          const active = path.startsWith(item.to);
          return (
            <div key={item.to}>
              <SidebarLists
                active={active}
                collapsed={collapsed}
                hoverOff={hoverOff}
                hoverOn={hoverOn}
                item={item}
              />
            </div>
          )
        })}
      </nav>

      {/* Divider */}
      <div
        className="mx-3 shrink-0"
        style={{ height: "1px", background: "var(--color-sidebar-border)" }}
      />

      

      {/* Logout */}
      <div className="px-2 py-3 flex flex-col gap-0.5">
        <button
          onClick={onLogout}
          className={[
            "flex items-center rounded-md w-full transition-colors",
            collapsed ? "justify-center p-2" : "gap-2.5 px-3 py-2 text-sm",
          ].join(" ")}
          style={{ color: "var(--color-sidebar-text)" }}
          onMouseEnter={(e) => hoverOn(e)}
          onMouseLeave={(e) => hoverOff(e)}
        >
          <FontAwesomeIcon icon={faRightFromBracket} />
          {!collapsed && <span>Logout</span>}
        </button>
      </div>
    </aside>
  );
}