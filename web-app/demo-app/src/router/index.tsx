import { Sidebar } from "@/features/layouts/Sidebar";
import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";
import Home from "./pages/Home";
import { WorkspacesPage } from "@/features/pages/Workspaces";

export default createBrowserRouter([
  {
    path: "/",
    element: (
      <div className="flex w-full h-full">
        <Sidebar />
        <Outlet />
      </div>
    ),
    children: [
      {
        index: true,
        element: <Navigate to="/home" replace />,
      },
      {
        path: "home",
        element: <Home />,
      },
      {
        path: "/workspaces",
        element: <WorkspacesPage />,
      }
    ],
  },
 
]);