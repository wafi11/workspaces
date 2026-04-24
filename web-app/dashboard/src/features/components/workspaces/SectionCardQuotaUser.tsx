import { useProfileQuota } from "@/features/api/auth";
import { Ico } from "@/features/components";
import type { PodVolumesType} from "@/types";
import type {ReactNode} from "react";

function pct(used: number, max: number) {
  if (max === 0) return 0;
  return Math.min(Math.round((used / max) * 100), 100);
}

function fmtRam(mb: number) {
  if (mb >= 1024) return `${(mb / 1024).toFixed(1)} GB`;
  return `${mb} MB`;
}

function Bar({ value }: { value: number }) {
  const color =
    value >= 90
      ? "#ef4444"
      : value >= 70
        ? "#f59e0b"
        : "var(--color-sidebar-accent)";
  return (
    <div
      className="w-full rounded-full overflow-hidden"
      style={{ height: 3, background: "var(--color-sidebar-border)" }}
    >
      <div
        className="h-full rounded-full transition-all duration-500"
        style={{ width: `${value}%`, background: color }}
      />
    </div>
  );
}

interface MetricProps {
  label: string;
  used: string;
  max: string;
  pct: number;
  icon: ReactNode;
}

function MetricCard({ label, used, max, pct: p, icon }: MetricProps) {
  return (
    <div
      className="flex flex-col gap-3 p-4 rounded-lg"
      style={{
        background: "var(--color-sidebar-bg)",
        border: "1px solid var(--color-sidebar-border)",
      }}
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span style={{ color: "var(--color-sidebar-text)" }}>{icon}</span>
          <span
            className="text-xs font-medium uppercase tracking-wider"
            style={{ color: "var(--color-sidebar-text)" }}
          >
            {label}
          </span>
        </div>
        <span
          className="text-xs tabular-nums"
          style={{ color: "var(--color-sidebar-text)" }}
        >
          {p}%
        </span>
      </div>
      <Bar value={p} />
      <div className="flex items-baseline justify-between">
        <span
          className="text-lg font-semibold tabular-nums tracking-tight"
          style={{ color: "var(--color-sidebar-text-active)" }}
        >
          {used}
        </span>
        <span
          className="text-xs tabular-nums"
          style={{ color: "var(--color-sidebar-text-muted)" }}
        >
          / {max}
        </span>
      </div>
    </div>
  );
}

interface SectionQuotaUserProps {
  volumes : PodVolumesType[]
}
export function SectionQuotaUser({ volumes }: SectionQuotaUserProps) {
  const { data: quota, isLoading } = useProfileQuota(); // tetap pake ini buat ws/cpu/ram

  if (isLoading) {
    return (
      <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-28 rounded-lg animate-pulse"
            style={{ background: "var(--color-sidebar-surface)" }} />
        ))}
      </div>
    );
  }

  if (!quota) return null;

  const q = quota;

  return (
    <div className="grid grid-cols-1 mt-4 p-4 sm:grid-cols-2 md:grid-cols-3 gap-3 lg:grid-cols-4">
      <MetricCard
        label="Workspaces"
        used={String(q.used_workspaces)}
        max={String(q.maxWorkspaces)}
        pct={pct(q.used_workspaces, q.maxWorkspaces)}
        icon={<Ico d="M3 7a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V7z" d2="M8 3v4M16 3v4" />}
      />
      <MetricCard
        label="CPU"
        used={`${q.used_cpu_cores} cores`}
        max={`${q.maxCpuCores} cores`}
        pct={pct(q.used_cpu_cores, q.maxCpuCores)}
        icon={<Ico d="M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2V9M9 21H5a2 2 0 0 1-2-2V9m0 0h18" />}
      />
      <MetricCard
        label="RAM"
        used={fmtRam(q.used_ram_mb)}
        max={fmtRam(q.maxRamMb)}
        pct={pct(q.used_ram_mb, q.maxRamMb)}
        icon={<Ico d="M6 19v-3m4 3v-7m4 7v-5m4 5V5M3 3h18v18H3z" />}
      />
      {
        volumes ? (
            <MetricCard
                label="Storage"
                used={fmtRam(volumes[0]?.used_mb)}
                max={fmtRam(volumes[0]?.capacity_mb / 10)}
                pct={pct(volumes[0]?.used_mb, volumes[0]?.capacity_mb)}
                icon={<Ico d="M22 12H2M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z" />}
            />
        ) : (
            <MetricCard
                label="Storage"
                used={q.used_storage_gb.toString()}
                max={q.maxStorageGb.toString()}
                pct={(q.used_storage_gb,q.maxStorageGb)}
                icon={<Ico d="M22 12H2M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z" />}
            />
        )
      }


    </div>
  );
}