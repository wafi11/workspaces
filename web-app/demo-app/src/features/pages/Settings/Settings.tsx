import { MainContainer } from "@/features/layouts/MainContainer";
import { SidebarSettings } from "@/features/layouts/SidebarSettings";
import { Outlet } from "react-router-dom";

export function Settings(){
    return (
        <>
            <SidebarSettings />
            <Outlet />
        </>
    )
}