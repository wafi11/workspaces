import type { NavItem } from "@/types";
import { Link } from "@tanstack/react-router";

interface SidebarListsProps {
    item: NavItem
    collapsed: boolean
    active: boolean
    hoverOn?: (e: React.MouseEvent<HTMLElement, MouseEvent>, active?: boolean) => void
    hoverOff?: (e: React.MouseEvent<HTMLElement, MouseEvent>, active?: boolean) => void
}

export function SidebarLists({ item, active, collapsed, hoverOff, hoverOn }: SidebarListsProps) {
    return (
        <Link
            to={item.to}
            className={[
                "flex items-center gap-2.5 mx-2 px-3 py-2 rounded-md text-sm",
                "transition-colors duration-100 relative",
                collapsed ? "justify-center" : "",
            ].join(" ")}
            style={{
                background: active ? "var(--color-sidebar-surface)" : "transparent",
                color: active ? "var(--color-sidebar-text-active)" : "var(--color-sidebar-text)",
            }}
            onMouseEnter={(e) => hoverOn?.(e, active)}
            onMouseLeave={(e) => hoverOff?.(e, active)}
        >
            {/* {active && (
                <span
                    className="absolute left-0 top-2 bottom-2 w-0.5 rounded-r-full"
                    style={{ background: "var(--color-sidebar-accent)" }}
                />
            )} */}
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
    )
}