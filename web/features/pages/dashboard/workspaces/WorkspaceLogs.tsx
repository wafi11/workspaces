import { cn } from "@/lib/utils";
import { RefObject } from "react";

interface WorkspaceLogsProps {
  logs: any[];
  logsEndRef: RefObject<HTMLDivElement | null>;
  logFilter: string;
}
export function WorkspaceLogs({ logs, logsEndRef }: WorkspaceLogsProps) {
  console.log("Rendering WorkspaceLogs with logs:", logs);
  return (
    <div className="flex-1 flex flex-col bg-[#050505]">
      {/* Toolbar tetap sama... */}

      <div className="flex-1 overflow-y-auto scrollbar-thin py-2">
        {logs.map((log, i) => {
          const time = log["@timestamp"]
            ? new Date(log["@timestamp"]).toLocaleTimeString([], {
                hour12: false,
              })
            : "--:--:--";

          const level = log.level || "INFO";
          const message = log.log || log.message || "";

          return (
            <div
              key={i}
              className="group flex gap-3 px-4 py-0.5 text-[11px] hover:bg-[#0f0f0f] font-mono border-l-2 border-transparent hover:border-blue-500/50 transition-all"
            >
              {/* Timestamp */}
              <span className="text-[#444] shrink-0 tabular-nums">
                [{time}]
              </span>

              {/* Level Badge */}
              <span
                className={cn(
                  "shrink-0 w-12 text-center rounded-[2px] text-[9px] font-bold leading-5 h-5 self-center tracking-tighter",
                  level === "INFO" &&
                    "text-blue-400 bg-blue-950/30 border border-blue-900/50",
                  level === "WARN" &&
                    "text-yellow-400 bg-yellow-950/30 border border-yellow-900/50",
                  level === "ERROR" &&
                    "text-red-400 bg-red-950/30 border border-red-900/50"
                )}
              >
                {level}
              </span>

              {/* Log Message */}
              <span className="text-[#ccc] break-all whitespace-pre-wrap leading-5">
                {message}
              </span>
            </div>
          );
        })}
        <div ref={logsEndRef} />
      </div>
    </div>
  );
}
