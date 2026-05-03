import { Sidebar } from "@/features/layouts/Sidebar";
import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";
import Home from "@/features/pages/Home";
import { WorkspacesPage } from "@/features/pages/Workspaces";
import { TemplatesPage } from "@/features/pages/Templates";
import { CreateWorkspacePage } from "@/features/pages/Workspaces/CreateWorkspacePage";
import { LoginPage } from "@/features/pages/Auth/Login";
import { RegisterPage } from "@/features/pages/Auth/Register";
import { Settings } from "@/features/pages/Settings/Settings";
import { SettingsProfile } from "@/features/pages/Settings/SettingsProfile";

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
      },
      {
        path : "/settings",
        element : <Settings />,
        children : [
          {
            path : "profile",
            element : <SettingsProfile />
          }
        ]
      }
    ],
  },
  {
    path : "/login",
    element : <LoginPage />
  },
  {
    path : "/register",
    element : <RegisterPage />
  }
]);