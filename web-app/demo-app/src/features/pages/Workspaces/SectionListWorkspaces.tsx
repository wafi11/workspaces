import type { Workspace } from "@/types";
import { faList, faTableCells } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useState } from "react";
import { GridView, ListView } from "./ListWorkspaces";
import { type ViewMode } from "./WorkspacesAction";

interface SectionListWorkspacesProps {
  workspaces: Workspace[];
}

export function SectionListWorkspaces({ workspaces }: SectionListWorkspacesProps) {
  const [view, setView] = useState<ViewMode>("list");

  return (
    <div
      className="rounded-lg m-4 overflow-hidden"
      style={{
        background: "#0a0a0a",
        border: "0.5px solid var(--color-app-border)",
      }}
    >
      {/* Header */}
      <div
        className="flex items-center justify-between px-4 py-2.5"
        style={{ borderBottom: "0.5px solid var(--color-app-border)" }}
      >
        <span className="text-xs" style={{ color: "var(--color-sidebar-text)" }}>
          {workspaces.length} workspace{workspaces.length !== 1 ? "s" : ""}
        </span>

        <div
          className="flex overflow-hidden rounded"
          style={{ border: "0.5px solid var(--color-app-border)" }}
        >
          {(["list", "grid"] as ViewMode[]).map((v) => (
            <button
              key={v}
              onClick={() => setView(v)}
              className="w-7 h-7 flex items-center justify-center"
              style={{
                background: view === v ? "var(--color-accent-surface)" : "var(--color-app-surface)",
                borderRight: v === "list" ? "0.5px solid var(--color-app-border)" : "none",
              }}
            >
              <FontAwesomeIcon
                icon={v === "list" ? faList : faTableCells}
                style={{
                  fontSize: "11px",
                  color: view === v ? "var(--color-accent-text)" : "var(--color-sidebar-text-muted)",
                }}
              />
            </button>
          ))}
        </div>
      </div>

      {workspaces.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <span className="text-xs" style={{ color: "var(--color-sidebar-text-muted)" }}>
            No workspaces found
          </span>
        </div>
      ) : view === "list" ? (
        <ListView workspaces={workspaces} />
      ) : (
        <GridView workspaces={workspaces} />
      )}
    </div>
  );
}