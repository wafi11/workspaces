"use client";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { useGetWorkspace } from "@/features/services/workspaces/api";
import { cn } from "@/lib/utils";
import { useEffect, useRef, useState } from "react";
import { InfoCard, InfoRow } from "./WorkspaceCard";
import { WorkspaceLogs } from "./WorkspaceLogs";
import { formatDate, relativeTime } from "@/utils/relativeTime";
import { API_URL } from "@/constants";

interface WorkspaceDetailsProps {
  slug: string;
}

type Tab = "terminal" | "logs";

export function WorkspaceDetails({ slug }: WorkspaceDetailsProps) {
  const { data } = useGetWorkspace(slug);
  const [activeTab, setActiveTab] = useState<Tab>("terminal");
  const [followLogs, setFollowLogs] = useState(true);
  const [logFilter, setLogFilter] = useState("");

  // State untuk menampung list logs
  const [logs, setLogs] = useState<any[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);

  const workspace = data?.data;
console.log(workspace)
  // --- LOGIC SSE ---
  useEffect(() => {
    if (!workspace?.namespace) return;

    const url = `${API_URL}logs/stream?namespace=${workspace.namespace}`;
    const eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      try {
        const logData = JSON.parse(event.data);
        console.log("Received log data:", logData);
        // Simpan objek utuh, bukan cuma string
        setLogs((prev) => [...prev.slice(-500), logData]);
      } catch (e) {
        console.error("Parse error", e);
      }
    };

    eventSource.onerror = (err) => {
      console.error("SSE Error:", err);
      eventSource.close();
    };

    // Cleanup: Tutup koneksi saat user pindah halaman atau komponen unmount
    return () => {
      eventSource.close();
    };
  }, [workspace?.name]); // Re-connect jika nama workspace berubah

  // --- AUTO SCROLL ---
  useEffect(() => {
    if (followLogs && activeTab === "logs") {
      logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [logs, followLogs, activeTab]);

  const filteredLogs = logs.filter((log) =>
    logFilter ? log.toLowerCase().includes(logFilter.toLowerCase()) : true
  );

  const statusColor =
    workspace?.status === "running"
      ? "text-green-400"
      : workspace?.status === "error"
        ? "text-red-400"
        : "text-yellow-400";

  return (
<>
    <TopBarAdmin
        title="Workspaces Details"
        description="Manage your workspace instance"
      />

      <div className="flex-1 overflow-y-auto p-5 grid grid-cols-2 gap-4 content-start">
        <InfoCard title="Workspace">
          <InfoRow label="ID" value={workspace?.id} mono />
          <InfoRow label="Name" value={workspace?.name} />
          <InfoRow
            label="Status"
            value={workspace?.status}
            className={statusColor}
          />
        </InfoCard>

        <InfoCard title="Timestamps">
          <InfoRow
            label="Created"
            value={formatDate(workspace?.created_at ?? "")}
          />
          <InfoRow
            label="Updated"
            value={relativeTime(workspace?.updated_at ?? "").toLocaleString()}
          />
          <InfoRow label="User ID" value={workspace?.user_id} mono />
        </InfoCard>

        <InfoCard title="Environment Variables" className="col-span-2">
          {workspace?.env_vars &&
          Object.entries(workspace.env_vars).length > 0 ? (
            Object.entries(workspace.env_vars).map(([k]) => (
              <InfoRow key={k} label={k} value="••••••••••" />
            ))
          ) : (
            <span className="text-gray-600 italic">
              No environment variables set
            </span>
          )}
        </InfoCard>
      </div>

      {/* Tabs Container */}
      <div className="flex flex-col h-[200px] px-5 border-t border-[#1a1a1a]">
        <div className="flex bg-[#0d0d0d] shrink-0">
          {(["terminal", "logs"] as Tab[]).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={cn(
                "px-4 py-2 text-[11px] uppercase tracking-widest border-b-2 transition-colors",
                activeTab === tab
                  ? "text-white border-blue-500"
                  : "text-[#555] border-transparent hover:text-[#999]"
              )}
            >
              {tab}
            </button>
          ))}
        </div>

        {/* Panel Content */}
        <div className="flex-1 bg-[#050505] overflow-hidden">
          {activeTab === "terminal" ? (
            workspace?.url && workspace.url !== "/" ? (
              <iframe
                src={workspace.url}
                className="w-full h-full border-0"
                title="workspace-terminal"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center flex-col gap-2 text-[#2a2a2a]">
                <span className="w-2 h-4 bg-blue-500 animate-pulse" />
                <span className="text-[11px]">terminal url not available</span>
              </div>
            )
          ) : (
            <WorkspaceLogs
              logFilter={logFilter}
              logs={filteredLogs}
              logsEndRef={logsEndRef}
            />
          )}
        </div>
      </div>
</>
  );
}
