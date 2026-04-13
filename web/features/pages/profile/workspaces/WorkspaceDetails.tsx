"use client";

import { WorkspaceDetails } from "../../dashboard/workspaces/WorkspaceDetails";

interface WorkspaceDetails {
  id: string;
}
export function WorkspaceDetailsUser({ id }: WorkspaceDetails) {
  return (
    <>
      <WorkspaceDetails slug={id} />
    </>
  );
}
