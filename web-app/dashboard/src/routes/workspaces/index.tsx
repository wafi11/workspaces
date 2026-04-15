import { TopbarAdmin } from "@/features/layout/TopbarDashboard";
import { SectionQuotaUser } from "@/features/pages/home/SectionCardQuotaUser";
import { SectionWorkspace } from "@/features/pages/home/SectionWorkspace";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/workspaces/")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
   <div className="flex flex-col gap-4 w-full">
        <TopbarAdmin title="Workspaces" className="py-3.75" />
        <main className="flex-1 flex flex-col min-w-0 overflow-y-auto p-4">
        <SectionQuotaUser />
        <SectionWorkspace />
      </main>
   </div>
  );
}
