import { SSE_URL } from "@/constants";
import { api } from "@/lib/api";
import { useEffect, useRef, useState } from "react";

export function DeployLog({ workspaceId }: { workspaceId?: string }) {
  const [logs, setLogs] = useState<string[]>([]);
  const bottomRef = useRef<HTMLDivElement>(null);
  const closeRef = useRef<(() => void) | null>(null);
  const lastConnectedId = useRef<string>("");

  useEffect(() => {
    if (!workspaceId) {
      setLogs([]);
      lastConnectedId.current = "";
    }
  }, [workspaceId]);

  useEffect(() => {
    if (!workspaceId || lastConnectedId.current === workspaceId) return;

    const startStream = async () => {
      lastConnectedId.current = workspaceId;
      try {
        const close = await api.sse(
          `${SSE_URL}/stream?namespace=ws-${workspaceId}`,
          (data: any) => {
            let rawLog = "";
            try {
              const parsed = typeof data === "string" ? JSON.parse(data) : data;
              rawLog = parsed.log || parsed.message || "";
            } catch {
              rawLog = typeof data === "string" ? data : "";
            }
            if (rawLog.includes(" F ")) {
              rawLog = rawLog.split(" F ")[1] || "";
            }

            const cleanLog = rawLog.trim();
            
            if (cleanLog) {
              setLogs((prev) => [...prev, cleanLog].slice(-200));
            }
          }
        );
        closeRef.current = close;
      } catch (err) {
        lastConnectedId.current = "";
      }
    };

    startStream();
    return () => {
      closeRef.current?.();
      closeRef.current = null;
    };
  }, [workspaceId]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <div className="flex flex-col h-full min-h-100 rounded-lg overflow-hidden bg-[#0d1117] shadow-xl">
      {/* Fake Terminal Header */}
      <div className="px-4 py-2 border-b border-white/10 flex justify-between items-center bg-[#161b22]">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-[#ff5f56]" />
          <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
          <div className="w-3 h-3 rounded-full bg-[#27c93f]" />
        </div>
        <span className="text-[10px] font-mono font-bold uppercase text-gray-500 tracking-tighter">
          Deployment Console — ws-{workspaceId?.slice(0, 8)}
        </span>
      </div>

      {/* Log View */}
      <div className="flex-1 overflow-y-auto p-4 font-mono text-[12px] selection:bg-blue-500/30">
        {!workspaceId ? (
          <div className="text-gray-600 italic">Standby - No active session</div>
        ) : (
          <div className="flex flex-col gap-0.5">
            {logs.map((l, i) => {
              // Cek apakah ini log error (mengandung 'N:' atau 'E:')
              const isError = l.includes("unrecognized option") || l.includes("error");
              
              return (
                <div key={i} className="flex gap-3 group">
                  {/* Line Number / Timestamp */}
                  <span className="text-gray-600 shrink-0 select-none w-14 text-right opacity-50 group-hover:opacity-100">
                    {new Date().toLocaleTimeString([], { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' })}
                  </span>
                  
                  {/* Log Message */}
                  <span className={`break-all ${isError ? 'text-red-400' : 'text-gray-300'}`}>
                    {l.split(' ').map((word, idx) => (
                      <span key={idx} className={word.startsWith('[') ? 'text-blue-400/80' : ''}>
                        {word}{' '}
                      </span>
                    ))}
                  </span>
                </div>
              );
            })}
            <div ref={bottomRef} className="h-4" />
          </div>
        )}
      </div>
    </div>
  );
}