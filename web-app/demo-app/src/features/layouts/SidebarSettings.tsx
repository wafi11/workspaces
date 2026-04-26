import { NavSettings } from "@/features/data";
import { SidebarLists } from "./SidebarList";
import { useLocation } from "react-router-dom";

export function SidebarSettings(){
      
    const location = useLocation()
    const path = location.pathname;
      
      
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
                "transition-[width] duration-200 ease-in-out w-56"
            ].join(" ")}
            style={{
                background: "var(--color-sidebar-bg)",
                borderRight: "1px solid var(--color-sidebar-border)",
            }}
            >
                <nav className="flex-1 py-2 flex flex-col gap-0.5 overflow-hidden">
                {NavSettings.map((item) => {
                const active = path.startsWith(item.to)

                const linkEl = (
                <SidebarLists active={active} collapsed={false} hoverOff={hoverOff} hoverOn={hoverOn} item={item}/>
                )

                return (
                    <div key={item.to}>{linkEl}</div>
                )
                })}
            </nav>
    </aside>
    )
}