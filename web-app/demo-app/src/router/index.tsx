import { Sidebar } from "@/features/layouts/Sidebar";
import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";
import Home from "@/features/pages/Home";
import { WorkspacesPage } from "@/features/pages/Workspaces";
import { TemplatesPage } from "@/features/pages/Templates";
import { CreateWorkspacePage } from "@/features/pages/Workspaces/CreateWorkspacePage";

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
        
      },
       {
        path: "/workspaces/create",
        element: <CreateWorkspacePage />,
      },
      {
        path: "/templates",
        element: <TemplatesPage />,
      }
    ],
  },
 
]);