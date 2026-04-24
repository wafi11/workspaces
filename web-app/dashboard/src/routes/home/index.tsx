import { useLogout, useProfile } from "@/features/api/auth";
import { SectionHeader } from "@/features/components/workspaces";
import { EmptyState } from "@/features/components/workspaces/EmptyState";
import { SectionQuotaUser } from "@/features/components/workspaces/SectionCardQuotaUser";
import { SectionTerminal } from "@/features/components/workspaces/SectionTerminal";
import { MainContainer } from "@/features/layout/MainContainer";
import { useGetMetricsStorageWorkspaces } from "@/hooks/useGetWorkspaceMatrics";
import { createFileRoute } from "@tanstack/react-router";
import type {PodVolumesType} from "@/types";

export const Route = createFileRoute("/home/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: profileData } = useProfile();
  const {metrics} = useGetMetricsStorageWorkspaces()

  return (
    
     
      <MainContainer>
        <SectionHeader />
        <SectionQuotaUser volumes={metrics?.volumes as PodVolumesType[]} />
        {/* <SectionWorkspace /> */}
        <SectionTerminal terminal_url={`https://${profileData?.terminal_url}`}/>
      </MainContainer>
  );
}
