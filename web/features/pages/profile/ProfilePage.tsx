"use client";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { useProfile, useProfileQuota } from "@/features/services/auth";
import { useGetWorkspacesUsers } from "@/features/services/workspaces/api";
import { getSystemGreeting } from "@/utils";
import { SectionTerminal } from "./SectionTerminal";
import { SectionUserQuota } from "./SectionUserQuota";
import { SectionWorkspaces } from "./SectionWorkspaces";
import { useWorkspaceSocket } from "@/hooks";
import { useQueryClient } from "@tanstack/react-query";

export function ProfilePage() {
  const { data: profile } = useProfile();
  const { data: dataQuota } = useProfileQuota();
  const { data: workspaceData } = useGetWorkspacesUsers();
  const greeting = getSystemGreeting();
  const queryClient = useQueryClient();

  useWorkspaceSocket((data) => {
    console.log(data);
    if (
      data.type === "workspace.running" ||
      data.type === "workspace.stopped"
    ) {
      queryClient.invalidateQueries({
        queryKey: ["workspace-users"],
      });
    }
  });

  return (
    <div className="container mx-auto pb-10">
      <TopBarAdmin
        description={`Welcome back, ${profile?.username || "User"}.\n ${greeting.sub}`}
        title={greeting.text}
      />
      <div className="space-y-8 mt-6">
        <SectionUserQuota data={dataQuota} />
        <SectionWorkspaces data={workspaceData} />
        <SectionTerminal url={profile?.terminal_url} />
      </div>
    </div>
  );
}
