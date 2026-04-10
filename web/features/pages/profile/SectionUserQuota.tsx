import { CardDetails } from "@/components/ui/card-details";
import { Database, Cpu, HardDrive, Layers, Activity, User } from "lucide-react";
import { UserQuota } from "@/types";

interface UserQuotaProps {
  data?: UserQuota;
}

export function SectionUserQuota({ data }: UserQuotaProps) {
  if (!data) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 animate-pulse">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="h-24 bg-muted rounded-md" />
        ))}
      </div>
    );
  }

  return (
    <section className="space-y-4 mt-4">
      <div className="flex items-center gap-2 mb-4">
        <Activity className="w-5 h-5 text-primary" />
        <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          Resource Quotas
        </h2>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* CPU Quota */}
        <CardDetails
          label="CPU Limit"
          value={`${data.used_cpu_cores}/${data.maxCpuCores} Cores`}
          sub="Maximum compute power"
          icon={Cpu}
        />

        {/* RAM Quota */}
        <CardDetails
          label="Memory"
          value={`${data.used_ram_mb / 1024}/${data.maxRamMb / 1024} GB`} // Konversi MB ke GB jika perlu
          sub={`${data.maxRamMb} MB total allocated`}
          icon={Database}
        />

        {/* Storage Quota */}
        <CardDetails
          label="Storage"
          value={`${data.used_storage_gb}/${data.maxStorageGb} GB`}
          sub="Persistent disk capacity"
          icon={HardDrive}
        />

        {/* Workspace Quota */}
        <CardDetails
          label="Workspaces"
          value={`${data.used_workspaces}/${data.maxWorkspaces}`}
          sub="Maximum active projects"
          icon={Layers}
        />
      </div>
    </section>
  );
}
